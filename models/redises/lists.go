package redises

import (
	"context"
	"gitlabee.4dogs.cn/common/redis"
)

type RedisList struct{}

//向链表的末尾插入一个元素
func (r *RedisList) SetListRPush(ctx context.Context, key string, value interface{}) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	if err := rediscli.RPush(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisList) SetListRPushAll(ctx context.Context, key string, value ...interface{}) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	if err := rediscli.RPush(ctx, key, value...).Err(); err != nil {
		return err
	}
	return nil
}

//向链表的第一位插入一个元素
func (r *RedisList) SetListLPush(ctx context.Context, key string, value interface{}) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	if err := rediscli.LPush(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

//删除链表的一个元素
func (r *RedisList) DeleteListLRem(ctx context.Context, key string, value interface{}) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	if err := rediscli.LRem(ctx, key, 0, value).Err(); err != nil {
		return err
	}
	return nil
}

//从链表取出元素
func (r *RedisList) GetListLRange(ctx context.Context, key string, auditWaitList *[]string) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	data := rediscli.LRange(ctx, key, 0, -1)
	if data == nil {
		return nil
	}
	*auditWaitList = data.Val()
	return nil
}

//从链表移除元素
func (r *RedisList) ListLRem(ctx context.Context, key string, nums int, value interface{}) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	if err := rediscli.LRem(ctx, key, int64(nums), value).Err(); err == nil {
		return nil
	}
	return nil
}

//从链表移除元素
func (r *RedisList) DelKey(ctx context.Context, key string) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	if err := rediscli.Del(ctx, key).Err(); err == nil {
		return nil
	}
	return nil
}

//获取链表长度
func (r *RedisList) GetListLen(ctx context.Context, key string) (int64, error) {
	rediscli, err := redis.NewClient()
	if err != nil {
		return 0, err
	}
	getKeyLong := rediscli.LLen(ctx, key)
	if err := getKeyLong.Err(); err != nil {
		return 0, err
	}
	return getKeyLong.Val(), nil
}

//从队列取出
func (r *RedisList) GetListWaitApplyData(ctx context.Context, key string, limit int64, auditWaitList *[]string) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	var value []string
	for i := int64(0); i < limit; i++ {
		data := rediscli.LPop(ctx, key)
		if data.Val() == "" {
			break
		}
		value = append(value, data.Val())
	}
	*auditWaitList = value
	return nil
}

func (r *RedisList) GetListRangeAndLtrim(ctx context.Context, key string, start int64, stop int64, auditWaitList *[]string) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	data := rediscli.LRange(ctx, key, start, stop-1)
	if data == nil {
		return nil
	}
	*auditWaitList = data.Val()
	if err := rediscli.LTrim(ctx, key, stop, -1).Err(); err != nil {
		return err
	}
	return nil
}

//从队列中删除元素，然后添加元素
func (r *RedisList) AddAndDeleteListLRem(ctx context.Context, key string, delValue []int64, addvalue interface{}) error {
	rediscli, err := redis.NewClient()
	if err != nil {
		return err
	}
	//删除
	for i := 0; i < len(delValue); i++ {
		if err := rediscli.LRem(ctx, key, 0, delValue[i]).Err(); err != nil {
			return err
		}
	}
	//新增
	if err := rediscli.RPush(ctx, key, addvalue).Err(); err != nil {
		return err
	}
	return nil
}
