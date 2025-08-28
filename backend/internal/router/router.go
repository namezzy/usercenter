package router

import (
	"usercenter/internal/handler"
	"usercenter/internal/middleware"
	
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	// 创建处理器实例
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	adminHandler := handler.NewAdminHandler()
	
	// API版本组
	api := r.Group("/api/v1")
	{
		// 公开路由（不需要认证）
		public := api.Group("")
		{
			// 认证相关
			auth := public.Group("/auth")
			{
				auth.GET("/captcha", authHandler.GetCaptcha)
				auth.POST("/send-email-code", authHandler.SendEmailCode)
				auth.POST("/send-sms-code", authHandler.SendSMSCode)
				auth.POST("/register", authHandler.Register)
				auth.POST("/login", authHandler.Login)
			}
		}
		
		// 需要认证的路由
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// 认证相关（需要登录）
			auth := protected.Group("/auth")
			{
				auth.POST("/logout", authHandler.Logout)
				auth.POST("/refresh", authHandler.RefreshToken)
				auth.GET("/user", authHandler.GetUserInfo)
			}
			
			// 用户个人中心
			profile := protected.Group("/profile")
			{
				profile.GET("", userHandler.GetProfile)
				profile.PUT("", userHandler.UpdateProfile)
				profile.PUT("/password", userHandler.ChangePassword)
				profile.POST("/avatar", userHandler.UploadAvatar)
				profile.POST("/bind-email", userHandler.BindEmail)
				profile.POST("/bind-phone", userHandler.BindPhone)
				profile.GET("/devices", userHandler.GetDevices)
				profile.DELETE("/devices/:device_id", userHandler.RemoveDevice)
				profile.GET("/logs", userHandler.GetLogs)
			}
		}
		
		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.AdminMiddleware())
		{
			// 用户管理
			users := admin.Group("/users")
			{
				users.GET("", adminHandler.GetUsers)
				users.POST("", adminHandler.CreateUser)
				users.GET("/:id", adminHandler.GetUser)
				users.PUT("/:id", adminHandler.UpdateUser)
				users.DELETE("/:id", adminHandler.DeleteUser)
				users.PUT("/:id/status", adminHandler.UpdateUserStatus)
				users.PUT("/:id/reset-password", adminHandler.ResetUserPassword)
			}
			
			// 统计信息
			admin.GET("/statistics", adminHandler.GetStatistics)
		}
		
		// 超级管理员路由
		superAdmin := api.Group("/super-admin")
		superAdmin.Use(middleware.AuthMiddleware())
		superAdmin.Use(middleware.SuperAdminMiddleware())
		{
			// 系统管理功能
			// 这里可以添加只有超级管理员才能访问的功能
		}
	}
	
	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// 静态文件服务
	r.Static("/uploads", "./uploads")
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "usercenter",
		})
	})
}
