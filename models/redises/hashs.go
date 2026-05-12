package redises

import (
	"context"
	"gitlabee.4dogs.cn/common/redis"
)

type RedisHash struct{}

//新增hash类型数据操作

//向hash中插入一条key
func (r *RedisHash) SetHashHset(ctx context.Context, key string, value ...interface{}) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	if err := rediscli.HSet(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

//向hash中插入一条key
func (r *RedisHash) GetHashHGet(ctx context.Context, key string, field string) (string, error) {
	rediscli, err := redis.NewClient()
	if err != nil {
		return "", err
	}
	getVal := rediscli.HGet(ctx, key, field)
	if err := getVal.Err(); err != nil {
		return "", err
	}
	return getVal.Val(), err
}

//向hash中插入一条key
func (r *RedisHash) GetHashHGetAll(ctx context.Context, key string) (map[string]string, error) {
	rediscli, err := redis.NewClient()
	if err != nil {
		return nil, err
	}
	getVal := rediscli.HGetAll(ctx, key)
	if err := getVal.Err(); err != nil {
		return nil, err
	}
	return getVal.Val(), err
}

//向hash中插入一条key
func (r *RedisHash) SetHashHMset(ctx context.Context, key string, values ...interface{}) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	getVal := rediscli.HMSet(ctx, key, values)
	if err := getVal.Err(); err != nil {
		return err
	}
	return err
}

//向hash中删除一条key
func (r *RedisHash) HDel(ctx context.Context, key string, values string) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	getVal := rediscli.HDel(ctx, key, values)
	if err := getVal.Err(); err != nil {
		return err
	}
	return err
}
