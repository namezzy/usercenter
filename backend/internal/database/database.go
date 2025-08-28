package database

import (
	"fmt"
	"usercenter/internal/config"
	"usercenter/internal/models"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode, cfg.Timezone)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	DB = db
	
	// 自动迁移数据库
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	
	// 初始化基础数据
	if err := initBaseData(); err != nil {
		return fmt.Errorf("failed to init base data: %w", err)
	}
	
	return nil
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserDevice{},
		&models.UserLog{},
		&models.VerificationCode{},
		&models.SystemNotification{},
		&models.UserNotification{},
		&models.DataBackup{},
	)
}

func initBaseData() error {
	// 创建默认角色
	roles := []models.Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "拥有系统所有权限",
			Status:      1,
			Sort:        1,
		},
		{
			Name:        "管理员",
			Code:        "admin",
			Description: "拥有大部分管理权限",
			Status:      1,
			Sort:        2,
		},
		{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户权限",
			Status:      1,
			Sort:        3,
		},
	}
	
	for _, role := range roles {
		var existingRole models.Role
		if err := DB.Where("code = ?", role.Code).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := DB.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	
	// 创建默认权限
	permissions := []models.Permission{
		// 系统管理
		{Name: "系统管理", Code: "system", Type: "menu", Path: "/system", Sort: 1},
		{Name: "用户管理", Code: "system:user", Type: "menu", Path: "/system/user", Sort: 1},
		{Name: "角色管理", Code: "system:role", Type: "menu", Path: "/system/role", Sort: 2},
		{Name: "权限管理", Code: "system:permission", Type: "menu", Path: "/system/permission", Sort: 3},
		{Name: "系统日志", Code: "system:log", Type: "menu", Path: "/system/log", Sort: 4},
		
		// 用户管理权限
		{Name: "查看用户", Code: "system:user:view", Type: "button", Sort: 1},
		{Name: "新增用户", Code: "system:user:add", Type: "button", Sort: 2},
		{Name: "编辑用户", Code: "system:user:edit", Type: "button", Sort: 3},
		{Name: "删除用户", Code: "system:user:delete", Type: "button", Sort: 4},
		{Name: "重置密码", Code: "system:user:reset", Type: "button", Sort: 5},
		
		// 个人中心
		{Name: "个人中心", Code: "profile", Type: "menu", Path: "/profile", Sort: 2},
		{Name: "基本信息", Code: "profile:info", Type: "menu", Path: "/profile/info", Sort: 1},
		{Name: "修改密码", Code: "profile:password", Type: "menu", Path: "/profile/password", Sort: 2},
		{Name: "安全设置", Code: "profile:security", Type: "menu", Path: "/profile/security", Sort: 3},
	}
	
	for _, permission := range permissions {
		var existingPermission models.Permission
		if err := DB.Where("code = ?", permission.Code).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := DB.Create(&permission).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	
	return nil
}
