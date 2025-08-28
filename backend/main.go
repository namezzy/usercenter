package main

import (
	"log"
	"usercenter/internal/cache"
	"usercenter/internal/config"
	"usercenter/internal/database"
	"usercenter/internal/middleware"
	"usercenter/internal/router"
	
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title 用户中心API
// @version 1.0
// @description 用户中心系统API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// 初始化日志
	logger, err := initLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()
	
	// 初始化数据库
	if err := database.InitDB(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}
	
	// 初始化Redis
	if err := cache.InitRedis(&cfg.Redis); err != nil {
		logger.Fatal("Failed to init redis", zap.Error(err))
	}
	
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)
	
	// 创建Gin实例
	r := gin.New()
	
	// 添加中间件
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.HealthCheckMiddleware())
	r.Use(middleware.OperationLogMiddleware())
	r.Use(middleware.RateLimitMiddleware(cfg.Security.RateLimit.RequestsPerMinute, cfg.Security.RateLimit.Burst))
	
	// 设置路由
	router.SetupRoutes(r)
	
	// 404处理
	r.NoRoute(middleware.NoRouteMiddleware())
	
	// 启动服务器
	logger.Info("Starting server", zap.String("port", cfg.Server.Port))
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func initLogger(cfg *config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	
	if cfg.Server.Mode == "debug" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	
	return logger, err
}
