package redises

import (
	"context"
	"gitlabee.4dogs.cn/common/redis"
	"time"
)

type RedisCommon struct{}

//新增通用操作类型

// 查询key
func (r *RedisCommon) Keys(ctx context.Context, pattern string) ([]string, error) {
	rediscli, err := redis.NewClient()
	if err != nil {
		return nil, err
	}
	getKeys := rediscli.Keys(ctx, pattern)
	if err := getKeys.Err(); err != nil {
		return nil, err
	}
	return getKeys.Val(), nil
}

// 重命名key
func (r *RedisCommon) Rename(ctx context.Context, key string, newKey string) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	getKeys := rediscli.Rename(ctx, key, newKey)
	if err := getKeys.Err(); err != nil {
		return err
	}
	return nil
}

// 执行setnx
func (r *RedisCommon) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	rediscli, err := redis.NewClient()
	if err != nil {
		return false, err
	}
	if expiration == 0 {
		expiration = 20 * time.Second
	}
	setnx := rediscli.SetNX(ctx, key, value, expiration)
	if err := setnx.Err(); err != nil {
		return false, err
	}
	return setnx.Val(), nil
}

// 重命名key
func (r *RedisCommon) Del(ctx context.Context, keys ...string) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	getKeys := rediscli.Del(ctx, keys...)
	if err := getKeys.Err(); err != nil {
		return err
	}
	return nil
}

// 重命名key
func (r *RedisCommon) Expire(ctx context.Context, key string, timeout time.Duration) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	getKeys := rediscli.Expire(ctx, key, timeout)
	if err := getKeys.Err(); err != nil {
		return err
	}
	return nil
}

// Ping 检查Redis连接状态
func (r *RedisCommon) Ping(ctx context.Context) (string, error) {
	rediscli, err := redis.NewClient()
	if err != nil {
		return "", err
	}
	pong := rediscli.Ping(ctx)
	if err := pong.Err(); err != nil {
		return "", err
	}
	return pong.Val(), nil
}
