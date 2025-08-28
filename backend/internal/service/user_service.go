package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"usercenter/internal/config"
	"usercenter/internal/database"
	"usercenter/internal/models"
	"usercenter/pkg/captcha"
	"usercenter/pkg/crypto"
	
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
}

type UpdateProfileRequest struct {
	Nickname string     `json:"nickname"`
	Gender   int        `json:"gender"`
	Birthday *time.Time `json:"birthday"`
	Bio      string     `json:"bio"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type BindEmailRequest struct {
	Email     string `json:"email" binding:"required,email"`
	EmailCode string `json:"email_code" binding:"required"`
}

type BindPhoneRequest struct {
	Phone   string `json:"phone" binding:"required"`
	SMSCode string `json:"sms_code" binding:"required"`
}

type UserListQuery struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
	RoleCode string `form:"role_code"`
}

type UserListResponse struct {
	Total int64         `json:"total"`
	Items []models.User `json:"items"`
}

func NewUserService() *UserService {
	return &UserService{}
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Roles").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	
	// 清空敏感信息
	user.Password = ""
	
	return &user, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uuid.UUID, req *UpdateProfileRequest) error {
	var user models.User
	err := database.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return err
	}
	
	// 更新字段
	user.Nickname = req.Nickname
	user.Gender = req.Gender
	user.Birthday = req.Birthday
	user.Bio = req.Bio
	
	return database.DB.Save(&user).Error
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uuid.UUID, req *ChangePasswordRequest) error {
	var user models.User
	err := database.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return err
	}
	
	// 验证旧密码
	isValid, err := crypto.VerifyPassword(req.OldPassword, user.Password)
	if err != nil {
		return err
	}
	
	if !isValid {
		return errors.New("旧密码不正确")
	}
	
	// 加密新密码
	hashedPassword, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	
	// 更新密码
	user.Password = hashedPassword
	return database.DB.Save(&user).Error
}

// UploadAvatar 上传头像
func (s *UserService) UploadAvatar(userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	// 检查文件类型
	allowedTypes := config.GlobalConfig.Upload.AllowedTypes
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	
	allowed := false
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif"} {
		if fileExt == ext {
			allowed = true
			break
		}
	}
	
	if !allowed {
		return "", errors.New("不支持的文件类型")
	}
	
	// 检查文件大小（这里简化处理，实际应该解析config中的MaxSize）
	if file.Size > 10*1024*1024 { // 10MB
		return "", errors.New("文件大小超过限制")
	}
	
	// 创建上传目录
	uploadPath := config.GlobalConfig.Upload.Path
	avatarDir := filepath.Join(uploadPath, "avatars")
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return "", err
	}
	
	// 生成文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	filePath := filepath.Join(avatarDir, filename)
	
	// 保存文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	
	// 更新用户头像
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)
	err = database.DB.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error
	if err != nil {
		// 删除已上传的文件
		os.Remove(filePath)
		return "", err
	}
	
	return avatarURL, nil
}

// BindEmail 绑定邮箱
func (s *UserService) BindEmail(userID uuid.UUID, req *BindEmailRequest) error {
	// 验证邮箱验证码
	if !captcha.VerifyEmailCode(req.Email, req.EmailCode, "bind_email") {
		return errors.New("验证码错误或已过期")
	}
	
	// 检查邮箱是否已被使用
	var existingUser models.User
	err := database.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error
	if err == nil {
		return errors.New("该邮箱已被其他用户绑定")
	}
	
	// 更新用户邮箱
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"email":          req.Email,
		"email_verified": true,
	}).Error
}

// BindPhone 绑定手机号
func (s *UserService) BindPhone(userID uuid.UUID, req *BindPhoneRequest) error {
	// 验证短信验证码
	if !captcha.VerifySMSCode(req.Phone, req.SMSCode, "bind_phone") {
		return errors.New("验证码错误或已过期")
	}
	
	// 检查手机号是否已被使用
	var existingUser models.User
	err := database.DB.Where("phone = ? AND id != ?", req.Phone, userID).First(&existingUser).Error
	if err == nil {
		return errors.New("该手机号已被其他用户绑定")
	}
	
	// 更新用户手机号
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"phone":          req.Phone,
		"phone_verified": true,
	}).Error
}

// GetUserDevices 获取用户设备列表
func (s *UserService) GetUserDevices(userID uuid.UUID) ([]models.UserDevice, error) {
	var devices []models.UserDevice
	err := database.DB.Where("user_id = ?", userID).Order("last_active desc").Find(&devices).Error
	return devices, err
}

// RemoveUserDevice 移除用户设备
func (s *UserService) RemoveUserDevice(userID uuid.UUID, deviceID string) error {
	return database.DB.Where("user_id = ? AND device_id = ?", userID, deviceID).Delete(&models.UserDevice{}).Error
}

// GetUserLogs 获取用户操作日志
func (s *UserService) GetUserLogs(userID uuid.UUID, page, pageSize int) (*UserListResponse, error) {
	var logs []models.UserLog
	var total int64
	
	query := database.DB.Model(&models.UserLog{}).Where("user_id = ?", userID)
	
	// 获取总数
	query.Count(&total)
	
	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	
	// 转换为用户列表响应格式（这里需要适配）
	return &UserListResponse{
		Total: total,
		Items: []models.User{}, // 实际应该返回日志列表
	}, nil
}

// AdminGetUsers 管理员获取用户列表
func (s *UserService) AdminGetUsers(query *UserListQuery) (*UserListResponse, error) {
	var users []models.User
	var total int64
	
	db := database.DB.Model(&models.User{}).Preload("Roles")
	
	// 关键词搜索
	if query.Keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	
	// 状态筛选
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	
	// 角色筛选
	if query.RoleCode != "" {
		db = db.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("roles.code = ?", query.RoleCode)
	}
	
	// 获取总数
	db.Count(&total)
	
	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Order("created_at desc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	
	// 清除敏感信息
	for i := range users {
		users[i].Password = ""
	}
	
	return &UserListResponse{
		Total: total,
		Items: users,
	}, nil
}

// AdminCreateUser 管理员创建用户
func (s *UserService) AdminCreateUser(req *RegisterRequest) error {
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
		EmailVerified: req.Email != "",
		PhoneVerified: req.Phone != "",
	}
	
	if req.Nickname == "" {
		user.Nickname = req.Username
	}
	
	return database.DB.Create(&user).Error
}

// AdminUpdateUser 管理员更新用户
func (s *UserService) AdminUpdateUser(userID uuid.UUID, req *UpdateProfileRequest) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(req).Error
}

// AdminDeleteUser 管理员删除用户
func (s *UserService) AdminDeleteUser(userID uuid.UUID) error {
	return database.DB.Delete(&models.User{}, userID).Error
}

// AdminUpdateUserStatus 管理员更新用户状态
func (s *UserService) AdminUpdateUserStatus(userID uuid.UUID, status int) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Update("status", status).Error
}

// AdminResetUserPassword 管理员重置用户密码
func (s *UserService) AdminResetUserPassword(userID uuid.UUID, newPassword string) error {
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password":       hashedPassword,
		"login_attempts": 0,
		"locked_until":   nil,
		"status":         models.UserStatusNormal,
	}).Error
}
