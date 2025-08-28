package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
	
	"usercenter/internal/database"
	"usercenter/internal/models"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return ""
		},
		Output: io.Discard,
	})
}

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 读取请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		
		// 处理请求
		c.Next()
		
		// 记录操作日志
		go func() {
			duration := time.Since(start)
			
			// 获取用户信息
			var userID uuid.UUID
			if uid, exists := c.Get("user_id"); exists {
				userID = uid.(uuid.UUID)
			}
			
			// 构建请求详情
			details := map[string]interface{}{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"query":      c.Request.URL.RawQuery,
				"user_agent": c.Request.UserAgent(),
				"duration":   duration.String(),
				"status":     c.Writer.Status(),
			}
			
			// 如果是敏感操作，记录请求体
			if shouldLogRequestBody(c.Request.Method, c.Request.URL.Path) {
				if len(bodyBytes) > 0 && len(bodyBytes) < 1024 { // 限制大小
					var requestBody map[string]interface{}
					if err := json.Unmarshal(bodyBytes, &requestBody); err == nil {
						// 移除敏感字段
						delete(requestBody, "password")
						delete(requestBody, "old_password")
						delete(requestBody, "new_password")
						details["request_body"] = requestBody
					}
				}
			}
			
			detailsJSON, _ := json.Marshal(details)
			
			// 确定操作类型和状态
			action := getActionFromPath(c.Request.Method, c.Request.URL.Path)
			status := 1
			if c.Writer.Status() >= 400 {
				status = 2
			}
			
			// 保存日志
			userLog := models.UserLog{
				UserID:    userID,
				Action:    action,
				Module:    getModuleFromPath(c.Request.URL.Path),
				IP:        c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
				Details:   string(detailsJSON),
				Status:    status,
			}
			
			database.DB.Create(&userLog)
		}()
	}
}

// shouldLogRequestBody 判断是否应该记录请求体
func shouldLogRequestBody(method, path string) bool {
	// 只记录POST、PUT、PATCH请求的请求体
	if method != "POST" && method != "PUT" && method != "PATCH" {
		return false
	}
	
	// 敏感操作路径
	sensitivePaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/users",
		"/api/v1/profile",
		"/api/v1/admin/users",
	}
	
	for _, sensitivePath := range sensitivePaths {
		if len(path) >= len(sensitivePath) && path[:len(sensitivePath)] == sensitivePath {
			return true
		}
	}
	
	return false
}

// getActionFromPath 从路径获取操作类型
func getActionFromPath(method, path string) string {
	if path == "/api/v1/auth/login" {
		return "用户登录"
	}
	if path == "/api/v1/auth/logout" {
		return "用户登出"
	}
	if path == "/api/v1/auth/register" {
		return "用户注册"
	}
	if path == "/api/v1/profile/password" {
		return "修改密码"
	}
	if path == "/api/v1/profile/avatar" {
		return "更新头像"
	}
	
	switch method {
	case "GET":
		return "查看"
	case "POST":
		return "新增"
	case "PUT", "PATCH":
		return "修改"
	case "DELETE":
		return "删除"
	default:
		return "操作"
	}
}

// getModuleFromPath 从路径获取模块名称
func getModuleFromPath(path string) string {
	if len(path) >= 12 && path[:12] == "/api/v1/auth" {
		return "认证模块"
	}
	if len(path) >= 15 && path[:15] == "/api/v1/profile" {
		return "个人中心"
	}
	if len(path) >= 13 && path[:13] == "/api/v1/users" {
		return "用户管理"
	}
	if len(path) >= 13 && path[:13] == "/api/v1/admin" {
		return "系统管理"
	}
	if len(path) >= 13 && path[:13] == "/api/v1/roles" {
		return "角色管理"
	}
	if len(path) >= 18 && path[:18] == "/api/v1/permissions" {
		return "权限管理"
	}
	
	return "系统"
}
