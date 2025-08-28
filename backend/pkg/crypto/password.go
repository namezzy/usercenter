package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	
	"golang.org/x/crypto/argon2"
)

type Config struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

var DefaultConfig = &Config{
	Time:    1,
	Memory:  64 * 1024,
	Threads: 4,
	KeyLen:  32,
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	salt, err := generateRandomBytes(16)
	if err != nil {
		return "", err
	}
	
	hash := argon2.IDKey([]byte(password), salt, DefaultConfig.Time, DefaultConfig.Memory, DefaultConfig.Threads, DefaultConfig.KeyLen)
	
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", 
		argon2.Version, DefaultConfig.Memory, DefaultConfig.Time, DefaultConfig.Threads, b64Salt, b64Hash)
	
	return encodedHash, nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, encodedHash string) (bool, error) {
	config, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	
	otherHash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)
	
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

// generateRandomBytes 生成随机字节
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// decodeHash 解码hash
func decodeHash(encodedHash string) (c *Config, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, fmt.Errorf("invalid hash format")
	}
	
	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("incompatible version of argon2")
	}
	
	c = &Config{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Time, &c.Threads)
	if err != nil {
		return nil, nil, nil, err
	}
	
	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	c.KeyLen = uint32(len(salt))
	
	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	c.KeyLen = uint32(len(hash))
	
	return c, salt, hash, nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	bytes, err := generateRandomBytes(uint32(length))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
