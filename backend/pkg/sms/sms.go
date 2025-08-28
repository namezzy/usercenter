package sms

import (
	"fmt"
	
	"usercenter/internal/config"
	
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SMSService struct {
	client *sms.Client
	config *config.TencentSMSConfig
}

func NewSMSService(cfg *config.TencentSMSConfig) (*SMSService, error) {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)
	
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	
	client, err := sms.NewClient(credential, "ap-beijing", cpf)
	if err != nil {
		return nil, err
	}
	
	return &SMSService{
		client: client,
		config: cfg,
	}, nil
}

// SendVerificationCode 发送验证码短信
func (s *SMSService) SendVerificationCode(phone, code string) error {
	request := sms.NewSendSmsRequest()
	
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.SmsSdkAppId = common.StringPtr(s.config.AppID)
	request.SignName = common.StringPtr(s.config.SignName)
	request.TemplateId = common.StringPtr(s.config.TemplateID)
	request.TemplateParamSet = common.StringPtrs([]string{code, "5"}) // 验证码和有效期
	
	response, err := s.client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("SMS API error: %v", err)
	}
	if err != nil {
		return err
	}
	
	// 检查发送结果
	if len(response.Response.SendStatusSet) > 0 {
		status := response.Response.SendStatusSet[0]
		if *status.Code != "Ok" {
			return fmt.Errorf("SMS send failed: %s", *status.Message)
		}
	}
	
	return nil
}

// SendNotificationSMS 发送通知短信
func (s *SMSService) SendNotificationSMS(phone, message string) error {
	// 这里可以根据需要实现不同的短信模板
	// 目前简化处理，使用通用模板
	request := sms.NewSendSmsRequest()
	
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.SmsSdkAppId = common.StringPtr(s.config.AppID)
	request.SignName = common.StringPtr(s.config.SignName)
	// 这里需要配置通知类短信的模板ID
	request.TemplateId = common.StringPtr("notification_template_id")
	request.TemplateParamSet = common.StringPtrs([]string{message})
	
	response, err := s.client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("SMS API error: %v", err)
	}
	if err != nil {
		return err
	}
	
	// 检查发送结果
	if len(response.Response.SendStatusSet) > 0 {
		status := response.Response.SendStatusSet[0]
		if *status.Code != "Ok" {
			return fmt.Errorf("SMS send failed: %s", *status.Message)
		}
	}
	
	return nil
}

// ValidatePhoneNumber 验证手机号格式
func ValidatePhoneNumber(phone string) bool {
	// 简单的中国手机号验证
	if len(phone) != 11 {
		return false
	}
	
	// 检查是否以1开头
	if phone[0] != '1' {
		return false
	}
	
	// 检查第二位是否为有效数字
	validSecondDigits := []byte{'3', '4', '5', '6', '7', '8', '9'}
	secondDigit := phone[1]
	
	for _, digit := range validSecondDigits {
		if secondDigit == digit {
			return true
		}
	}
	
	return false
}

// FormatPhoneNumber 格式化手机号（添加+86前缀）
func FormatPhoneNumber(phone string) string {
	if len(phone) == 11 && phone[0] == '1' {
		return "+86" + phone
	}
	return phone
}
