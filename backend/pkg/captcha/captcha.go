package captcha

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
	
	"usercenter/internal/cache"
	"usercenter/internal/models"
	"usercenter/internal/database"
	
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

// ImageCaptcha 图形验证码配置
type ImageCaptcha struct {
	Width  int
	Height int
	Length int
}

// DefaultImageCaptcha 默认图形验证码配置
var DefaultImageCaptcha = &ImageCaptcha{
	Width:  240,
	Height: 80,
	Length: 4,
}

// GenerateImageCaptcha 生成图形验证码
func GenerateImageCaptcha() (id, b64s string, err error) {
	driver := base64Captcha.NewDriverDigit(
		DefaultImageCaptcha.Height,
		DefaultImageCaptcha.Width,
		DefaultImageCaptcha.Length,
		0.7,
		80,
	)
	
	c := base64Captcha.NewCaptcha(driver, store)
	return c.Generate()
}

// VerifyImageCaptcha 验证图形验证码
func VerifyImageCaptcha(id, answer string) bool {
	return store.Verify(id, answer, true)
}

// GenerateEmailCode 生成邮箱验证码
func GenerateEmailCode(email, purpose string) (string, error) {
	code := generateNumericCode(6)
	
	// 保存到数据库
	verificationCode := models.VerificationCode{
		Type:      "email",
		Target:    email,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(15 * time.Minute), // 15分钟过期
	}
	
	if err := database.DB.Create(&verificationCode).Error; err != nil {
		return "", err
	}
	
	// 保存到Redis缓存
	key := fmt.Sprintf("email_code:%s:%s", email, purpose)
	if err := cache.Set(key, code, 15*time.Minute); err != nil {
		return "", err
	}
	
	return code, nil
}

// GenerateSMSCode 生成短信验证码
func GenerateSMSCode(phone, purpose string) (string, error) {
	code := generateNumericCode(6)
	
	// 保存到数据库
	verificationCode := models.VerificationCode{
		Type:      "sms",
		Target:    phone,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟过期
	}
	
	if err := database.DB.Create(&verificationCode).Error; err != nil {
		return "", err
	}
	
	// 保存到Redis缓存
	key := fmt.Sprintf("sms_code:%s:%s", phone, purpose)
	if err := cache.Set(key, code, 5*time.Minute); err != nil {
		return "", err
	}
	
	return code, nil
}

// VerifyEmailCode 验证邮箱验证码
func VerifyEmailCode(email, code, purpose string) bool {
	// 从Redis验证
	key := fmt.Sprintf("email_code:%s:%s", email, purpose)
	storedCode, err := cache.Get(key)
	if err != nil {
		return false
	}
	
	if storedCode != code {
		return false
	}
	
	// 验证成功后删除验证码
	cache.Del(key)
	
	// 更新数据库记录为已使用
	database.DB.Model(&models.VerificationCode{}).
		Where("type = ? AND target = ? AND code = ? AND purpose = ? AND used = false", "email", email, code, purpose).
		Update("used", true)
	
	return true
}

// VerifySMSCode 验证短信验证码
func VerifySMSCode(phone, code, purpose string) bool {
	// 从Redis验证
	key := fmt.Sprintf("sms_code:%s:%s", phone, purpose)
	storedCode, err := cache.Get(key)
	if err != nil {
		return false
	}
	
	if storedCode != code {
		return false
	}
	
	// 验证成功后删除验证码
	cache.Del(key)
	
	// 更新数据库记录为已使用
	database.DB.Model(&models.VerificationCode{}).
		Where("type = ? AND target = ? AND code = ? AND purpose = ? AND used = false", "sms", phone, code, purpose).
		Update("used", true)
	
	return true
}

// CheckCodeSendFrequency 检查验证码发送频率
func CheckCodeSendFrequency(target, codeType string) (bool, time.Duration) {
	key := fmt.Sprintf("send_frequency:%s:%s", codeType, target)
	
	exists, err := cache.Exists(key)
	if err != nil {
		return true, 0
	}
	
	if exists {
		// 获取剩余时间
		ttl := cache.RDB.TTL(cache.RDB.Context(), key).Val()
		return false, ttl
	}
	
	// 设置发送频率限制（60秒）
	cache.Set(key, "1", 60*time.Second)
	return true, 0
}

// generateNumericCode 生成数字验证码
func generateNumericCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10))
	}
	return code
}

// GenerateBase64Captcha 生成Base64编码的验证码
func GenerateBase64Captcha() (string, string, error) {
	id, b64s, err := GenerateImageCaptcha()
	if err != nil {
		return "", "", err
	}
	
	// 将Base64字符串转换为完整的Data URL
	dataURL := "data:image/png;base64," + b64s
	
	return id, dataURL, nil
}
