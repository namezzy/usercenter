package email

import (
	"fmt"
	"strconv"
	
	"usercenter/internal/config"
	
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.SMTPConfig
}

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

func NewEmailService(cfg *config.SMTPConfig) *EmailService {
	return &EmailService{config: cfg}
}

// SendEmail 发送邮件
func (s *EmailService) SendEmail(msg *EmailMessage) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.Username)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	
	if msg.IsHTML {
		m.SetBody("text/html", msg.Body)
	} else {
		m.SetBody("text/plain", msg.Body)
	}
	
	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	
	return d.DialAndSend(m)
}

// SendVerificationCode 发送验证码邮件
func (s *EmailService) SendVerificationCode(email, code, purpose string) error {
	subject := getEmailSubject(purpose)
	body := getEmailBody(code, purpose)
	
	msg := &EmailMessage{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}
	
	return s.SendEmail(msg)
}

// SendWelcomeEmail 发送欢迎邮件
func (s *EmailService) SendWelcomeEmail(email, username string) error {
	subject := "欢迎注册用户中心"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>欢迎注册用户中心</h2>
			<p>亲爱的 %s，</p>
			<p>欢迎您注册我们的用户中心！您的账号已成功创建。</p>
			<p>现在您可以登录并开始使用我们的服务。</p>
			<p>如果您有任何问题，请随时联系我们的客服团队。</p>
			<br>
			<p>祝好！</p>
			<p>用户中心团队</p>
		</body>
		</html>
	`, username)
	
	msg := &EmailMessage{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}
	
	return s.SendEmail(msg)
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *EmailService) SendPasswordResetEmail(email, resetLink string) error {
	subject := "密码重置请求"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>密码重置请求</h2>
			<p>您好，</p>
			<p>我们收到了您的密码重置请求。请点击下面的链接来重置您的密码：</p>
			<p><a href="%s" style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">重置密码</a></p>
			<p>如果您无法点击上面的按钮，请复制以下链接到浏览器地址栏：</p>
			<p>%s</p>
			<p>此链接将在30分钟后失效。</p>
			<p>如果您没有请求密码重置，请忽略此邮件。</p>
			<br>
			<p>祝好！</p>
			<p>用户中心团队</p>
		</body>
		</html>
	`, resetLink, resetLink)
	
	msg := &EmailMessage{
		To:      []string{email},
		Subject: subject,
		Body:    body,
		IsHTML:  true,
	}
	
	return s.SendEmail(msg)
}

// SendNotificationEmail 发送通知邮件
func (s *EmailService) SendNotificationEmail(email, title, content string) error {
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>%s</h2>
			<div>%s</div>
			<br>
			<p>祝好！</p>
			<p>用户中心团队</p>
		</body>
		</html>
	`, title, content)
	
	msg := &EmailMessage{
		To:      []string{email},
		Subject: title,
		Body:    body,
		IsHTML:  true,
	}
	
	return s.SendEmail(msg)
}

func getEmailSubject(purpose string) string {
	switch purpose {
	case "register":
		return "用户注册验证码"
	case "reset_password":
		return "密码重置验证码"
	case "login":
		return "登录验证码"
	case "bind_email":
		return "邮箱绑定验证码"
	default:
		return "验证码"
	}
}

func getEmailBody(code, purpose string) string {
	action := getActionText(purpose)
	
	return fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #007bff; text-align: center;">验证码</h2>
				<p>您好，</p>
				<p>您正在进行%s操作，验证码为：</p>
				<div style="text-align: center; margin: 30px 0;">
					<span style="font-size: 36px; font-weight: bold; color: #007bff; letter-spacing: 5px; border: 2px dashed #007bff; padding: 15px 25px; display: inline-block;">%s</span>
				</div>
				<p style="color: #666;">验证码有效期为15分钟，请及时使用。</p>
				<p style="color: #666;">如果这不是您本人操作，请忽略此邮件。</p>
				<br>
				<p>祝好！</p>
				<p>用户中心团队</p>
				<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
				<p style="font-size: 12px; color: #999; text-align: center;">
					此邮件由系统自动发送，请勿直接回复。
				</p>
			</div>
		</body>
		</html>
	`, action, code)
}

func getActionText(purpose string) string {
	switch purpose {
	case "register":
		return "用户注册"
	case "reset_password":
		return "密码重置"
	case "login":
		return "登录验证"
	case "bind_email":
		return "邮箱绑定"
	default:
		return "身份验证"
	}
}
