package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// User 用户模型
type User struct {
	BaseModel
	Username        string     `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50"`
	Email           string     `json:"email" gorm:"uniqueIndex" validate:"email"`
	Phone           string     `json:"phone" gorm:"uniqueIndex"`
	Password        string     `json:"-" gorm:"not null" validate:"required,min=8"`
	Nickname        string     `json:"nickname" gorm:"size:100"`
	Avatar          string     `json:"avatar"`
	Gender          int        `json:"gender" gorm:"default:0"` // 0:未知 1:男 2:女
	Birthday        *time.Time `json:"birthday"`
	Bio             string     `json:"bio" gorm:"size:500"`
	Status          int        `json:"status" gorm:"default:1"` // 1:正常 2:禁用 3:锁定
	EmailVerified   bool       `json:"email_verified" gorm:"default:false"`
	PhoneVerified   bool       `json:"phone_verified" gorm:"default:false"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	LastLoginIP     string     `json:"last_login_ip"`
	LoginAttempts   int        `json:"login_attempts" gorm:"default:0"`
	LockedUntil     *time.Time `json:"locked_until"`
	TwoFactorSecret string     `json:"-"`
	TwoFactorEnabled bool      `json:"two_factor_enabled" gorm:"default:false"`
	
	// 关联关系
	Roles       []Role       `json:"roles" gorm:"many2many:user_roles;"`
	UserDevices []UserDevice `json:"user_devices"`
	UserLogs    []UserLog    `json:"user_logs"`
}

// Role 角色模型
type Role struct {
	BaseModel
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Code        string       `json:"code" gorm:"uniqueIndex;not null"`
	Description string       `json:"description"`
	Status      int          `json:"status" gorm:"default:1"`
	Sort        int          `json:"sort" gorm:"default:0"`
	
	// 关联关系
	Users       []User       `json:"users" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

// Permission 权限模型
type Permission struct {
	BaseModel
	Name        string `json:"name" gorm:"not null"`
	Code        string `json:"code" gorm:"uniqueIndex;not null"`
	Type        string `json:"type" gorm:"not null"` // menu, button, data
	ParentID    *uuid.UUID `json:"parent_id"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort" gorm:"default:0"`
	Status      int    `json:"status" gorm:"default:1"`
	
	// 关联关系
	Parent   *Permission `json:"parent" gorm:"foreignKey:ParentID"`
	Children []Permission `json:"children" gorm:"foreignKey:ParentID"`
	Roles    []Role      `json:"roles" gorm:"many2many:role_permissions;"`
}

// UserDevice 用户设备模型
type UserDevice struct {
	BaseModel
	UserID      uuid.UUID `json:"user_id" gorm:"not null"`
	DeviceID    string    `json:"device_id" gorm:"not null"`
	DeviceType  string    `json:"device_type"` // web, mobile, app
	DeviceName  string    `json:"device_name"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	LastActive  time.Time `json:"last_active"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	
	// 关联关系
	User User `json:"user"`
}

// UserLog 用户日志模型
type UserLog struct {
	BaseModel
	UserID     uuid.UUID `json:"user_id"`
	Action     string    `json:"action" gorm:"not null"`
	Module     string    `json:"module"`
	IP         string    `json:"ip"`
	UserAgent  string    `json:"user_agent"`
	Details    string    `json:"details" gorm:"type:text"`
	Status     int       `json:"status"` // 1:成功 2:失败
	
	// 关联关系
	User User `json:"user"`
}

// VerificationCode 验证码模型
type VerificationCode struct {
	BaseModel
	Type      string    `json:"type" gorm:"not null"` // email, sms, captcha
	Target    string    `json:"target" gorm:"not null"` // 邮箱或手机号
	Code      string    `json:"code" gorm:"not null"`
	Purpose   string    `json:"purpose" gorm:"not null"` // register, reset_password, login
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Used      bool      `json:"used" gorm:"default:false"`
	IP        string    `json:"ip"`
}

// SystemNotification 系统通知模型
type SystemNotification struct {
	BaseModel
	Title    string    `json:"title" gorm:"not null"`
	Content  string    `json:"content" gorm:"type:text"`
	Type     string    `json:"type" gorm:"not null"` // info, warning, error, success
	Priority int       `json:"priority" gorm:"default:1"` // 1:低 2:中 3:高
	Status   int       `json:"status" gorm:"default:1"` // 1:启用 2:禁用
	
	// 发送设置
	SendToAll   bool      `json:"send_to_all" gorm:"default:false"`
	SendAt      *time.Time `json:"send_at"`
	ExpireAt    *time.Time `json:"expire_at"`
	
	// 关联关系
	Recipients []UserNotification `json:"recipients"`
}

// UserNotification 用户通知关联模型
type UserNotification struct {
	BaseModel
	UserID         uuid.UUID `json:"user_id" gorm:"not null"`
	NotificationID uuid.UUID `json:"notification_id" gorm:"not null"`
	ReadAt         *time.Time `json:"read_at"`
	IsRead         bool      `json:"is_read" gorm:"default:false"`
	
	// 关联关系
	User         User               `json:"user"`
	Notification SystemNotification `json:"notification"`
}

// DataBackup 数据备份模型
type DataBackup struct {
	BaseModel
	Name        string    `json:"name" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // full, incremental
	FilePath    string    `json:"file_path" gorm:"not null"`
	FileSize    int64     `json:"file_size"`
	Status      int       `json:"status" gorm:"default:1"` // 1:成功 2:失败 3:进行中
	StartedAt   time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Error       string    `json:"error"`
	CreatedBy   uuid.UUID `json:"created_by"`
}

// 用户状态常量
const (
	UserStatusNormal   = 1
	UserStatusDisabled = 2
	UserStatusLocked   = 3
)

// 性别常量
const (
	GenderUnknown = 0
	GenderMale    = 1
	GenderFemale  = 2
)

// 权限类型常量
const (
	PermissionTypeMenu   = "menu"
	PermissionTypeButton = "button"
	PermissionTypeData   = "data"
)

// 通知类型常量
const (
	NotificationTypeInfo    = "info"
	NotificationTypeWarning = "warning"
	NotificationTypeError   = "error"
	NotificationTypeSuccess = "success"
)
