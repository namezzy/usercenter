package handler

import (
	"net/http"
	"strconv"
	
	"usercenter/internal/middleware"
	"usercenter/internal/service"
	
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// GetProfile 获取个人资料
// @Summary 获取个人资料
// @Description 获取当前用户的详细资料
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "用户资料"
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户资料失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"message": "获取用户资料成功",
	})
}

// UpdateProfile 更新个人资料
// @Summary 更新个人资料
// @Description 更新当前用户的基本信息
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.UpdateProfileRequest true "更新信息"
// @Success 200 {object} map[string]interface{} "更新结果"
// @Router /profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新用户资料失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新用户资料成功",
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的登录密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.ChangePasswordRequest true "密码信息"
// @Success 200 {object} map[string]interface{} "修改结果"
// @Router /profile/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err := h.userService.ChangePassword(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码修改成功",
	})
}

// UploadAvatar 上传头像
// @Summary 上传头像
// @Description 上传用户头像
// @Tags 用户
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param avatar formData file true "头像文件"
// @Success 200 {object} map[string]interface{} "上传结果"
// @Router /profile/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请选择头像文件",
		})
		return
	}
	
	avatarURL, err := h.userService.UploadAvatar(userID, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"avatar_url": avatarURL,
		},
		"message": "头像上传成功",
	})
}

// BindEmail 绑定邮箱
// @Summary 绑定邮箱
// @Description 绑定用户邮箱
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.BindEmailRequest true "邮箱信息"
// @Success 200 {object} map[string]interface{} "绑定结果"
// @Router /profile/bind-email [post]
func (h *UserHandler) BindEmail(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	var req service.BindEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err := h.userService.BindEmail(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "邮箱绑定成功",
	})
}

// BindPhone 绑定手机号
// @Summary 绑定手机号
// @Description 绑定用户手机号
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.BindPhoneRequest true "手机号信息"
// @Success 200 {object} map[string]interface{} "绑定结果"
// @Router /profile/bind-phone [post]
func (h *UserHandler) BindPhone(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	var req service.BindPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err := h.userService.BindPhone(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "手机号绑定成功",
	})
}

// GetDevices 获取设备列表
// @Summary 获取设备列表
// @Description 获取当前用户的登录设备列表
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "设备列表"
// @Router /profile/devices [get]
func (h *UserHandler) GetDevices(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	devices, err := h.userService.GetUserDevices(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取设备列表失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": devices,
		"message": "获取设备列表成功",
	})
}

// RemoveDevice 移除设备
// @Summary 移除设备
// @Description 移除指定的登录设备
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param device_id path string true "设备ID"
// @Success 200 {object} map[string]interface{} "移除结果"
// @Router /profile/devices/{device_id} [delete]
func (h *UserHandler) RemoveDevice(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "设备ID不能为空",
		})
		return
	}
	
	err := h.userService.RemoveUserDevice(userID, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "移除设备失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "设备移除成功",
	})
}

// GetLogs 获取操作日志
// @Summary 获取操作日志
// @Description 获取当前用户的操作日志
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{} "操作日志"
// @Router /profile/logs [get]
func (h *UserHandler) GetLogs(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	logs, err := h.userService.GetUserLogs(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取操作日志失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": logs,
		"message": "获取操作日志成功",
	})
}
