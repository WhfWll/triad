package enums

const (
	// 平台token - 防止缓存击穿 lock
	RedisPlatformTokenLock = "platform_token_lock:"
	// 平台token - 缓存
	RedisPlatformToken = "platform_token"

	// 用户最新操作时间锁，防止频繁修改
	RedisUserUpdateOperateTimeLock = "last_operate_time_lock"

	// 在线用户名单，用来处理退出后的用户 Authorization 的校验
	RedisUserAlive = "redis_user_alive"
)
