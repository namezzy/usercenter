package middleware

import (
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流器结构
type RateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(requestsPerMinute int, burst int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), burst),
	}
}

// RateLimitMiddleware 全局限流中间件
func RateLimitMiddleware(requestsPerMinute int, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute, burst)
	
	return func(c *gin.Context) {
		if !limiter.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimitMiddleware IP限流中间件
func IPRateLimitMiddleware(requestsPerMinute int, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), burst)
			limiters[ip] = limiter
		}
		
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// UserRateLimitMiddleware 用户限流中间件
func UserRateLimitMiddleware(requestsPerMinute int, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			// 如果没有用户ID，使用IP限流
			ip := c.ClientIP()
			limiter, exists := limiters[ip]
			if !exists {
				limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), burst)
				limiters[ip] = limiter
			}
			
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code":    429,
					"message": "请求过于频繁，请稍后再试",
				})
				c.Abort()
				return
			}
		} else {
			userIDStr := userID.(string)
			limiter, exists := limiters[userIDStr]
			if !exists {
				limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), burst)
				limiters[userIDStr] = limiter
			}
			
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code":    429,
					"message": "请求过于频繁，请稍后再试",
				})
				c.Abort()
				return
			}
		}
		
		c.Next()
	}
}
