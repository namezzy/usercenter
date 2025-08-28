package handler

import (
	"net/http"
	"strconv"
	
	"usercenter/internal/middleware"
	"usercenter/internal/service"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	userService *service.UserService
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		userService: service.NewUserService(),
	}
}

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 管理员获取用户列表
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param status query int false "用户状态"
// @Param role_code query string false "角色代码"
// @Success 200 {object} map[string]interface{} "用户列表"
// @Router /admin/users [get]
func (h *AdminHandler) GetUsers(c *gin.Context) {
	var query service.UserListQuery
	
	// 绑定查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}
	
	result, err := h.userService.AdminGetUsers(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取用户列表失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
		"message": "获取用户列表成功",
	})
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 管理员创建新用户
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.RegisterRequest true "用户信息"
// @Success 200 {object} map[string]interface{} "创建结果"
// @Router /admin/users [post]
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err := h.userService.AdminCreateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户创建成功",
	})
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 管理员获取指定用户的详细信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Success 200 {object} map[string]interface{} "用户详情"
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID格式错误",
		})
		return
	}
	
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"message": "获取用户详情成功",
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 管理员更新用户信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Param request body service.UpdateProfileRequest true "更新信息"
// @Success 200 {object} map[string]interface{} "更新结果"
// @Router /admin/users/{id} [put]
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID格式错误",
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
	
	err = h.userService.AdminUpdateUser(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新用户失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户更新成功",
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 管理员删除用户
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Success 200 {object} map[string]interface{} "删除结果"
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	// 检查是否是超级管理员
	userRole, _ := middleware.GetUserRole(c)
	if userRole != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "只有超级管理员可以删除用户",
		})
		return
	}
	
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID格式错误",
		})
		return
	}
	
	// 不允许删除自己
	currentUserID, _ := middleware.GetUserID(c)
	if currentUserID == userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不能删除自己",
		})
		return
	}
	
	err = h.userService.AdminDeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除用户失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户删除成功",
	})
}

// UpdateUserStatus 更新用户状态
// @Summary 更新用户状态
// @Description 管理员更新用户状态（启用/禁用/锁定）
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Param request body map[string]int true "状态信息"
// @Success 200 {object} map[string]interface{} "更新结果"
// @Router /admin/users/{id}/status [put]
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID格式错误",
		})
		return
	}
	
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	// 不允许修改自己的状态
	currentUserID, _ := middleware.GetUserID(c)
	if currentUserID == userID {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不能修改自己的状态",
		})
		return
	}
	
	err = h.userService.AdminUpdateUserStatus(userID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新用户状态失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户状态更新成功",
	})
}

// ResetUserPassword 重置用户密码
// @Summary 重置用户密码
// @Description 管理员重置用户密码
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "用户ID"
// @Param request body map[string]string true "新密码"
// @Success 200 {object} map[string]interface{} "重置结果"
// @Router /admin/users/{id}/reset-password [put]
func (h *AdminHandler) ResetUserPassword(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户ID格式错误",
		})
		return
	}
	
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	err = h.userService.AdminResetUserPassword(userID, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "重置密码失败",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码重置成功",
	})
}

// GetStatistics 获取统计信息
// @Summary 获取统计信息
// @Description 获取用户中心的统计信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "统计信息"
// @Router /admin/statistics [get]
func (h *AdminHandler) GetStatistics(c *gin.Context) {
	// 这里可以实现获取各种统计信息
	// 如用户总数、今日新增、今日活跃等
	
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total_users":  0,
			"today_new":    0,
			"today_active": 0,
			"online_users": 0,
		},
		"message": "获取统计信息成功",
	})
}
