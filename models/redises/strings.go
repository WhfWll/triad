package redises

import (
	"context"
	"gitlabee.4dogs.cn/common/redis"
	"time"
)

type RedisStr struct{}

//新增string类型数据
//redis.DB为redis单库，如果是主从模式使用redis.Cluster

func (r *RedisStr) SetString1(ctx context.Context, key string, value string, expirattime int) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	err = rediscli.Set(ctx, key, value, time.Duration(expirattime)*time.Second).Err()
	return err
}

func (r *RedisStr) SetString2(ctx context.Context, key string, value string, expirattime int) error {
	rediscli, err := redis.NewClient("local")
	if err != nil {
		return err
	}
	err = rediscli.Set(ctx, key, value, time.Duration(expirattime)*time.Second).Err()
	return err
}

//获取值
func (r *RedisStr) Get(ctx context.Context, key string) (string, error) {
	rediscli, err := redis.NewClient()
	if err != nil {
		return "", err
	}
	getRedisString := rediscli.Get(ctx, key)
	if err := rediscli.Get(ctx, key).Err(); err != nil {
		return "", err
	}
	return getRedisString.Val(), nil
}
