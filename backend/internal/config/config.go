package config

import (
	"time"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	SMTP     SMTPConfig     `mapstructure:"smtp"`
	SMS      SMSConfig      `mapstructure:"sms"`
	Security SecurityConfig `mapstructure:"security"`
	Upload   UploadConfig   `mapstructure:"upload"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret  string        `mapstructure:"secret"`
	Expires time.Duration `mapstructure:"expires"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type SMSConfig struct {
	Tencent TencentSMSConfig `mapstructure:"tencent"`
}

type TencentSMSConfig struct {
	SecretID   string `mapstructure:"secret_id"`
	SecretKey  string `mapstructure:"secret_key"`
	AppID      string `mapstructure:"app_id"`
	SignName   string `mapstructure:"sign_name"`
	TemplateID string `mapstructure:"template_id"`
}

type SecurityConfig struct {
	MaxLoginAttempts   int           `mapstructure:"max_login_attempts"`
	LockDuration       time.Duration `mapstructure:"lock_duration"`
	PasswordMinLength  int           `mapstructure:"password_min_length"`
	RateLimit          RateLimitConfig `mapstructure:"rate_limit"`
}

type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	Burst             int `mapstructure:"burst"`
}

type UploadConfig struct {
	MaxSize      string   `mapstructure:"max_size"`
	AllowedTypes []string `mapstructure:"allowed_types"`
	Path         string   `mapstructure:"path"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

var GlobalConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	
	// 设置默认值
	setDefaults()
	
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	
	GlobalConfig = &config
	return &config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.timezone", "UTC")
	
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	
	viper.SetDefault("jwt.expires", "24h")
	
	viper.SetDefault("security.max_login_attempts", 5)
	viper.SetDefault("security.lock_duration", "30m")
	viper.SetDefault("security.password_min_length", 8)
	viper.SetDefault("security.rate_limit.requests_per_minute", 60)
	viper.SetDefault("security.rate_limit.burst", 10)
	
	viper.SetDefault("upload.max_size", "10MB")
	viper.SetDefault("upload.path", "./uploads")
	
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.filename", "./logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.max_backups", 10)
}
