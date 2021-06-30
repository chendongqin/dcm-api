package redis_value

import (
	"dongchamao/global/logger"
	redisHelp "dongchamao/global/redis"
	"dongchamao/global/utils"
	"github.com/gomodule/redigo/redis"
)

func getRedis() *redisHelp.RClient {
	return redisHelp.Redis("default")
}

// 判断是否存在此 Key
func KeyBool(err error) bool {
	if err == redis.ErrNil {
		return false
	}
	return true
}

type LockInRedis struct {
	Sec string
	Key string
	Ex  int
}

func NewRedisLock(key string, ex int) *LockInRedis {
	return &LockInRedis{
		Key: key,
		Sec: utils.GetRandomStringNew(16),
		Ex:  ex,
	}
}

func (l *LockInRedis) Lock() (bool, error) {
	return getRedis().Set(l.Key, l.Sec, "NX", "EX", l.Ex)
}

func (l *LockInRedis) Release() {
	sec, err := getRedis().GetString(l.Key)
	if KeyBool(err) {
		if sec == l.Sec {
			getRedis().Del(l.Key)
		}
	} else {
		logger.Error("can not find redis lock key:", l.Key)
	}
}
