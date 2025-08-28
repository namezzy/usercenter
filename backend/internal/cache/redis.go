package cache

import (
	"context"
	"fmt"
	"time"
	
	"usercenter/internal/config"
	
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var ctx = context.Background()

func InitRedis(cfg *config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	
	// 测试连接
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	
	return nil
}

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// Del 删除缓存
func Del(key string) error {
	return RDB.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	result, err := RDB.Exists(ctx, key).Result()
	return result > 0, err
}

// Incr 递增
func Incr(key string) (int64, error) {
	return RDB.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return RDB.Expire(ctx, key, expiration).Err()
}

// SetNX 如果键不存在则设置
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return RDB.SetNX(ctx, key, value, expiration).Result()
}

// HSet 哈希设置
func HSet(key, field string, value interface{}) error {
	return RDB.HSet(ctx, key, field, value).Err()
}

// HGet 哈希获取
func HGet(key, field string) (string, error) {
	return RDB.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func HGetAll(key string) (map[string]string, error) {
	return RDB.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) error {
	return RDB.HDel(ctx, key, fields...).Err()
}

// ZAdd 有序集合添加
func ZAdd(key string, score float64, member interface{}) error {
	return RDB.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// ZRangeByScore 按分数范围获取有序集合
func ZRangeByScore(key string, min, max string) ([]string, error) {
	return RDB.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

// LPush 列表左推
func LPush(key string, values ...interface{}) error {
	return RDB.LPush(ctx, key, values...).Err()
}

// RPop 列表右弹
func RPop(key string) (string, error) {
	return RDB.RPop(ctx, key).Result()
}

// LLen 列表长度
func LLen(key string) (int64, error) {
	return RDB.LLen(ctx, key).Result()
}

// SAdd 集合添加
func SAdd(key string, members ...interface{}) error {
	return RDB.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(key string) ([]string, error) {
	return RDB.SMembers(ctx, key).Result()
}

// SRem 集合删除
func SRem(key string, members ...interface{}) error {
	return RDB.SRem(ctx, key, members...).Err()
}
