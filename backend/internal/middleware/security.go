package middleware

import (
	"net/http"
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"}, // 允许的源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// SecurityHeadersMiddleware 安全头中间件
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")
		
		// 防止MIME类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")
		
		// XSS保护
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// 强制HTTPS（在生产环境中启用）
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// 引用策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "服务器内部错误",
		})
	})
}

// NoRouteMiddleware 404处理中间件
func NoRouteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "接口不存在",
		})
	}
}

// HealthCheckMiddleware 健康检查中间件
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"time":   time.Now().Unix(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
