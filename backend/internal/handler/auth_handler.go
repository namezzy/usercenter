package handler

import (
	"net/http"
	
	"usercenter/internal/middleware"
	"usercenter/internal/service"
	"usercenter/pkg/captcha"
	
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// GetCaptcha 获取图形验证码
// @Summary 获取图形验证码
// @Description 获取用于登录/注册的图形验证码
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "验证码信息"
// @Router /auth/captcha [get]
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	id, b64s, err := captcha.GenerateBase64Captcha()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成验证码失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"captcha_id":  id,
			"captcha_img": b64s,
		},
		"message": "获取验证码成功",
	})
}

// SendEmailCode 发送邮箱验证码
// @Summary 发送邮箱验证码
// @Description 发送邮箱验证码用于注册、绑定等操作
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body map[string]string true "邮箱和用途"
// @Success 200 {object} map[string]interface{} "发送结果"
// @Router /auth/send-email-code [post]
func (h *AuthHandler) SendEmailCode(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required,email"`
		Purpose string `json:"purpose" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err := h.authService.SendEmailCode(req.Email, req.Purpose)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "验证码发送成功",
	})
}

// SendSMSCode 发送短信验证码
// @Summary 发送短信验证码
// @Description 发送短信验证码用于注册、绑定等操作
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body map[string]string true "手机号和用途"
// @Success 200 {object} map[string]interface{} "发送结果"
// @Router /auth/send-sms-code [post]
func (h *AuthHandler) SendSMSCode(c *gin.Context) {
	var req struct {
		Phone   string `json:"phone" binding:"required"`
		Purpose string `json:"purpose" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err := h.authService.SendSMSCode(req.Phone, req.Purpose)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "验证码发送成功",
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} map[string]interface{} "注册结果"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "登录结果"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	// 设置设备信息
	req.DeviceInfo.IP = c.ClientIP()
	req.DeviceInfo.UserAgent = c.GetHeader("User-Agent")
	
	if req.DeviceInfo.DeviceID == "" {
		req.DeviceInfo.DeviceID = c.GetHeader("X-Device-ID")
	}
	if req.DeviceInfo.DeviceType == "" {
		req.DeviceInfo.DeviceType = "web"
	}
	
	resp, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": resp,
		"message": "登录成功",
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出接口
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "登出结果"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token, _ := middleware.GetToken(c)
	deviceID, _ := middleware.GetDeviceID(c)
	
	err := h.authService.Logout(token, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "登出失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登出成功",
	})
}

// RefreshToken 刷新Token
// @Summary 刷新Token
// @Description 刷新访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "刷新结果"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	token, _ := middleware.GetToken(c)
	
	newToken, err := h.authService.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"token": newToken,
		},
		"message": "Token刷新成功",
	})
}

// GetUserInfo 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的基本信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "用户信息"
// @Router /auth/user [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	userService := service.NewUserService()
	user, err := userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"message": "获取用户信息成功",
	})
}
