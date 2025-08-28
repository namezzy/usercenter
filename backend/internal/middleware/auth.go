package middleware

import (
	"net/http"
	"strings"
	"time"
	
	"usercenter/internal/cache"
	"usercenter/pkg/jwt"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请先登录",
			})
			c.Abort()
			return
		}
		
		// 移除Bearer前缀
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}
		
		// 解析Token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token无效",
			})
			c.Abort()
			return
		}
		
		// 检查Token是否在黑名单中
		blacklistKey := "token_blacklist:" + token
		exists, _ := cache.Exists(blacklistKey)
		if exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token已失效",
			})
			c.Abort()
			return
		}
		
		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("device_id", claims.DeviceID)
		c.Set("token", token)
		
		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			if strings.HasPrefix(token, "Bearer ") {
				token = token[7:]
			}
			
			claims, err := jwt.ParseToken(token)
			if err == nil {
				// 检查Token是否在黑名单中
				blacklistKey := "token_blacklist:" + token
				exists, _ := cache.Exists(blacklistKey)
				if !exists {
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("role", claims.Role)
					c.Set("device_id", claims.DeviceID)
					c.Set("token", token)
				}
			}
		}
		
		c.Next()
	}
}

// RoleMiddleware 角色权限中间件
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足",
			})
			c.Abort()
			return
		}
		
		roleStr := userRole.(string)
		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}
		
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "权限不足",
		})
		c.Abort()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("super_admin", "admin")
}

// SuperAdminMiddleware 超级管理员权限中间件
func SuperAdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("super_admin")
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	
	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	
	name, ok := username.(string)
	return name, ok
}

// GetUserRole 从上下文中获取用户角色
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	
	roleStr, ok := role.(string)
	return roleStr, ok
}

// GetDeviceID 从上下文中获取设备ID
func GetDeviceID(c *gin.Context) (string, bool) {
	deviceID, exists := c.Get("device_id")
	if !exists {
		return "", false
	}
	
	id, ok := deviceID.(string)
	return id, ok
}

// GetToken 从上下文中获取Token
func GetToken(c *gin.Context) (string, bool) {
	token, exists := c.Get("token")
	if !exists {
		return "", false
	}
	
	tokenStr, ok := token.(string)
	return tokenStr, ok
}

// BlacklistToken 将Token加入黑名单
func BlacklistToken(token string, expiration time.Duration) error {
	blacklistKey := "token_blacklist:" + token
	return cache.Set(blacklistKey, "1", expiration)
}
