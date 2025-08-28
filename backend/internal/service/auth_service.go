package service

import (
	"errors"
	"fmt"
	"time"
	
	"usercenter/internal/cache"
	"usercenter/internal/config"
	"usercenter/internal/database"
	"usercenter/internal/models"
	"usercenter/pkg/captcha"
	"usercenter/pkg/crypto"
	"usercenter/pkg/email"
	"usercenter/pkg/jwt"
	"usercenter/pkg/sms"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	emailService *email.EmailService
	smsService   *sms.SMSService
}

type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
	DeviceInfo  DeviceInfo `json:"device_info"`
}

type RegisterRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Email        string `json:"email" binding:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password" binding:"required,min=8"`
	Nickname     string `json:"nickname"`
	EmailCode    string `json:"email_code"`
	SMSCode      string `json:"sms_code"`
	CaptchaID    string `json:"captcha_id" binding:"required"`
	CaptchaCode  string `json:"captcha_code" binding:"required"`
}

type DeviceInfo struct {
	DeviceID   string `json:"device_id"`
	DeviceType string `json:"device_type"` // web, mobile, app
	DeviceName string `json:"device_name"`
	UserAgent  string `json:"user_agent"`
	IP         string `json:"ip"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	ExpiresAt int64       `json:"expires_at"`
}

func NewAuthService() *AuthService {
	cfg := config.GlobalConfig
	emailSvc := email.NewEmailService(&cfg.SMTP)
	smsSvc, _ := sms.NewSMSService(&cfg.SMS.Tencent)
	
	return &AuthService{
		emailService: emailSvc,
		smsService:   smsSvc,
	}
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 验证图形验证码
	if req.CaptchaID != "" && req.CaptchaCode != "" {
		if !captcha.VerifyImageCaptcha(req.CaptchaID, req.CaptchaCode) {
			return nil, errors.New("验证码错误")
		}
	}
	
	// 查找用户
	var user models.User
	err := database.DB.Preload("Roles").Where("username = ? OR email = ? OR phone = ?", 
		req.Username, req.Username, req.Username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}
	
	// 检查账号状态
	if user.Status == models.UserStatusDisabled {
		return nil, errors.New("账号已被禁用")
	}
	
	// 检查账号是否被锁定
	if user.Status == models.UserStatusLocked && user.LockedUntil != nil && 
		time.Now().Before(*user.LockedUntil) {
		return nil, fmt.Errorf("账号已被锁定，请在 %s 后重试", user.LockedUntil.Format("2006-01-02 15:04:05"))
	}
	
	// 验证密码
	isValid, err := crypto.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, err
	}
	
	if !isValid {
		// 记录登录失败
		s.recordLoginFailure(&user, req.DeviceInfo.IP)
		return nil, errors.New("用户名或密码错误")
	}
	
	// 重置登录失败次数
	s.resetLoginFailures(&user)
	
	// 更新最后登录信息
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = req.DeviceInfo.IP
	database.DB.Save(&user)
	
	// 记录设备信息
	s.recordDeviceInfo(&user, req.DeviceInfo)
	
	// 获取用户角色
	var roleCode string
	if len(user.Roles) > 0 {
		roleCode = user.Roles[0].Code
	}
	
	// 生成Token
	token, err := jwt.GenerateToken(user.ID, user.Username, roleCode, req.DeviceInfo.DeviceID)
	if err != nil {
		return nil, err
	}
	
	// 计算过期时间
	expiresAt := time.Now().Add(config.GlobalConfig.JWT.Expires).Unix()
	
	return &LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) error {
	// 验证图形验证码
	if !captcha.VerifyImageCaptcha(req.CaptchaID, req.CaptchaCode) {
		return errors.New("验证码错误")
	}
	
	// 验证邮箱验证码（如果提供了邮箱）
	if req.Email != "" {
		if req.EmailCode == "" {
			return errors.New("请输入邮箱验证码")
		}
		if !captcha.VerifyEmailCode(req.Email, req.EmailCode, "register") {
			return errors.New("邮箱验证码错误或已过期")
		}
	}
	
	// 验证短信验证码（如果提供了手机号）
	if req.Phone != "" {
		if req.SMSCode == "" {
			return errors.New("请输入短信验证码")
		}
		if !captcha.VerifySMSCode(req.Phone, req.SMSCode, "register") {
			return errors.New("短信验证码错误或已过期")
		}
	}
	
	// 检查用户名是否已存在
	var existingUser models.User
	err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return errors.New("用户名已存在")
	}
	
	// 检查邮箱是否已存在
	if req.Email != "" {
		err = database.DB.Where("email = ?", req.Email).First(&existingUser).Error
		if err == nil {
			return errors.New("邮箱已存在")
		}
	}
	
	// 检查手机号是否已存在
	if req.Phone != "" {
		err = database.DB.Where("phone = ?", req.Phone).First(&existingUser).Error
		if err == nil {
			return errors.New("手机号已存在")
		}
	}
	
	// 加密密码
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return err
	}
	
	// 创建用户
	user := models.User{
		Username:      req.Username,
		Email:         req.Email,
		Phone:         req.Phone,
		Password:      hashedPassword,
		Nickname:      req.Nickname,
		Status:        models.UserStatusNormal,
		EmailVerified: req.Email != "" && req.EmailCode != "",
		PhoneVerified: req.Phone != "" && req.SMSCode != "",
	}
	
	if req.Nickname == "" {
		user.Nickname = req.Username
	}
	
	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	// 创建用户
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 分配默认角色（普通用户）
	var userRole models.Role
	if err := tx.Where("code = ?", "user").First(&userRole).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	if err := tx.Model(&user).Association("Roles").Append(&userRole); err != nil {
		tx.Rollback()
		return err
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	
	// 发送欢迎邮件
	if req.Email != "" {
		go func() {
			s.emailService.SendWelcomeEmail(req.Email, user.Username)
		}()
	}
	
	return nil
}

// SendEmailCode 发送邮箱验证码
func (s *AuthService) SendEmailCode(email, purpose string) error {
	// 检查发送频率
	canSend, remaining := captcha.CheckCodeSendFrequency(email, "email")
	if !canSend {
		return fmt.Errorf("发送过于频繁，请在 %d 秒后重试", int(remaining.Seconds()))
	}
	
	// 生成验证码
	code, err := captcha.GenerateEmailCode(email, purpose)
	if err != nil {
		return err
	}
	
	// 发送邮件
	return s.emailService.SendVerificationCode(email, code, purpose)
}

// SendSMSCode 发送短信验证码
func (s *AuthService) SendSMSCode(phone, purpose string) error {
	// 验证手机号格式
	if !sms.ValidatePhoneNumber(phone) {
		return errors.New("手机号格式不正确")
	}
	
	// 检查发送频率
	canSend, remaining := captcha.CheckCodeSendFrequency(phone, "sms")
	if !canSend {
		return fmt.Errorf("发送过于频繁，请在 %d 秒后重试", int(remaining.Seconds()))
	}
	
	// 生成验证码
	code, err := captcha.GenerateSMSCode(phone, purpose)
	if err != nil {
		return err
	}
	
	// 发送短信
	return s.smsService.SendVerificationCode(phone, code)
}

// Logout 用户登出
func (s *AuthService) Logout(token string, deviceID string) error {
	// 将Token加入黑名单
	expiration := config.GlobalConfig.JWT.Expires
	err := cache.Set("token_blacklist:"+token, "1", expiration)
	if err != nil {
		return err
	}
	
	// 更新设备状态
	if deviceID != "" {
		database.DB.Model(&models.UserDevice{}).
			Where("device_id = ?", deviceID).
			Update("is_active", false)
	}
	
	return nil
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(token string) (string, error) {
	// 解析当前Token
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return "", err
	}
	
	// 检查Token是否在黑名单中
	blacklistKey := "token_blacklist:" + token
	exists, _ := cache.Exists(blacklistKey)
	if exists {
		return "", errors.New("Token已失效")
	}
	
	// 生成新Token
	newToken, err := jwt.GenerateToken(claims.UserID, claims.Username, claims.Role, claims.DeviceID)
	if err != nil {
		return "", err
	}
	
	// 将旧Token加入黑名单
	expiration := config.GlobalConfig.JWT.Expires
	cache.Set(blacklistKey, "1", expiration)
	
	return newToken, nil
}

// recordLoginFailure 记录登录失败
func (s *AuthService) recordLoginFailure(user *models.User, ip string) {
	user.LoginAttempts++
	
	// 如果失败次数达到限制，锁定账号
	maxAttempts := config.GlobalConfig.Security.MaxLoginAttempts
	if user.LoginAttempts >= maxAttempts {
		lockDuration := config.GlobalConfig.Security.LockDuration
		lockedUntil := time.Now().Add(lockDuration)
		user.LockedUntil = &lockedUntil
		user.Status = models.UserStatusLocked
	}
	
	database.DB.Save(user)
}

// resetLoginFailures 重置登录失败次数
func (s *AuthService) resetLoginFailures(user *models.User) {
	user.LoginAttempts = 0
	user.LockedUntil = nil
	if user.Status == models.UserStatusLocked {
		user.Status = models.UserStatusNormal
	}
	database.DB.Save(user)
}

// recordDeviceInfo 记录设备信息
func (s *AuthService) recordDeviceInfo(user *models.User, deviceInfo DeviceInfo) {
	// 查找或创建设备记录
	var device models.UserDevice
	err := database.DB.Where("user_id = ? AND device_id = ?", user.ID, deviceInfo.DeviceID).First(&device).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新设备记录
		device = models.UserDevice{
			UserID:     user.ID,
			DeviceID:   deviceInfo.DeviceID,
			DeviceType: deviceInfo.DeviceType,
			DeviceName: deviceInfo.DeviceName,
			IP:         deviceInfo.IP,
			UserAgent:  deviceInfo.UserAgent,
			LastActive: time.Now(),
			IsActive:   true,
		}
		database.DB.Create(&device)
	} else {
		// 更新现有设备记录
		device.IP = deviceInfo.IP
		device.UserAgent = deviceInfo.UserAgent
		device.LastActive = time.Now()
		device.IsActive = true
		database.DB.Save(&device)
	}
}
