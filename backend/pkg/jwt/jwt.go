package jwt

import (
	"errors"
	"time"
	
	"usercenter/internal/config"
	
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	DeviceID string    `json:"device_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uuid.UUID, username, role, deviceID string) (string, error) {
	cfg := config.GlobalConfig
	if cfg == nil {
		return "", errors.New("config not initialized")
	}
	
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.Expires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "usercenter",
			Subject:   userID.String(),
			ID:        uuid.New().String(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GlobalConfig
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

// RefreshToken 刷新Token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	
	// 检查Token是否即将过期（剩余时间小于30分钟）
	if time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		return "", errors.New("token does not need refresh")
	}
	
	return GenerateToken(claims.UserID, claims.Username, claims.Role, claims.DeviceID)
}

// ValidateToken 验证Token是否有效
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}
