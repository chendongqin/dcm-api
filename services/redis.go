package services

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	json2 "encoding/json"
	"fmt"
	"github.com/astaxie/beego/cache"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisService struct {
}

func NewRedisService() *RedisService {
	return new(RedisService)
}

func (this *RedisService) GetRedis() redis.Conn {
	redisConn := global.Cache.GetInstance().(redis.Conn)
	return redisConn
}

func (this *RedisService) Hset(key string, field string, value string) error {
	redisConn := this.GetRedis()
	defer redisConn.Close()
	_, err := redisConn.Do("hset", key, field, value)
	return err
}
func (this *RedisService) HsetWithTimeout(key string, field string, value string, timeout int) (err error) {
	redisConn := this.GetRedis()
	defer redisConn.Close()
	_, err = redisConn.Do("hset", key, field, value)
	if err != nil {
		return
	}
	_, err = redisConn.Do("EXPIRE", key, timeout)
	return
}

func (this *RedisService) Hget(key string, field string) string {
	redisConn := this.GetRedis()
	defer redisConn.Close()
	var str string
	if v, err := redisConn.Do("hget", key, field); err == nil {
		if v != nil {
			str, _ = redis.String(v, nil)
		} else {
			str = ""
		}
		return str
	}
	return ""
}

func (this *RedisService) Incr(key string) (reply interface{}, err error) {
	redisConn := this.GetRedis()
	defer redisConn.Close()
	return redisConn.Do("incr", key)
}

type ZSetItem struct {
	Score float64
	Data  interface{}
}

type ZSetItems []ZSetItem

func (this *RedisService) SAdd(key string, value string) (int, error) {
	conn := this.GetRedis()
	defer conn.Close()
	reply, err := conn.Do("sadd", key, value)
	if err != nil {
		return 0, err
	}
	return cache.GetInt(reply), err
}

func (this *RedisService) SCard(key string) (int, error) {
	conn := this.GetRedis()
	defer conn.Close()
	reply, err := conn.Do("scard", key)
	if err != nil {
		return 0, err
	}
	return cache.GetInt(reply), err
}

func (this *RedisService) SPop(key string, count int) (string, error) {
	conn := this.GetRedis()
	defer conn.Close()
	reply, err := conn.Do("spop", key)
	if err != nil {
		return "", err
	}
	return cache.GetString(reply), err
}

func (this *RedisService) Rename(key string, newKey string) (err error) {
	conn := this.GetRedis()
	defer conn.Close()
	_, err = conn.Do("rename", key, newKey)
	if err != nil {
		return
	}
	return
}

func (this *RedisService) ZAdd(key string, kv ZSetItems, timeout int) (err error) {
	conn := this.GetRedis()
	defer conn.Close()
	//args := make([]interface{}, 0)
	//args = append(args, key)
	for _, v := range kv {
		var val string
		if str, ok := v.Data.(string); ok {
			val = str
		} else {
			var json []byte
			json, err = json2.Marshal(v.Data)
			if err != nil {
				return
			}
			val = string(json)
		}
		//args = append(args, i)
		//args = append(args)
		err = conn.Send("ZADD", key, v.Score, val)
		if err != nil {
			return
		}
	}
	err = conn.Send("EXPIRE", key, timeout)
	if err != nil {
		return
	}
	err = conn.Flush()
	return
}

func (this *RedisService) ZCard(key string) (num int, err error) {
	conn := this.GetRedis()
	defer conn.Close()
	reply, err := conn.Do("ZCARD", key)
	if err != nil {
		return
	}
	return cache.GetInt(reply), err
}

func (this *RedisService) ZRevRange(key string, start int, end int, isWithScore bool) (reply interface{}, err error) {
	conn := this.GetRedis()
	defer conn.Close()
	if isWithScore {
		reply, err = conn.Do("zrevrange", key, start, end, "WITHSCORES")
	} else {
		reply, err = conn.Do("zrevrange", key, start, end)
	}
	return
}

func (this *RedisService) ZRevRangeByScore(key string, min float64, max float64, isWithScore bool, limit ...int) (reply interface{}, err error) {
	conn := this.GetRedis()
	defer conn.Close()
	args := []interface{}{
		key,
		max,
		min,
	}
	if isWithScore {
		args = append(args, "WITHSCORES")
	}
	if len(limit) >= 2 {
		args = append(args, "LIMIT", limit[0], limit[1])
	}
	reply, err = conn.Do("zrevrangebyscore", args...)
	return
}

func (this *RedisService) ZRangeByScore(key string, min float64, max float64, isWithScore bool, limit ...int) (reply interface{}, err error) {
	conn := this.GetRedis()
	defer conn.Close()
	args := []interface{}{
		key,
		min,
		max,
	}
	if isWithScore {
		args = append(args, "WITHSCORES")
	}
	if len(limit) >= 2 {
		args = append(args, "LIMIT", limit[0], limit[1])
	}
	reply, err = conn.Do("zrangebyscore", args...)
	return
}

func (this *RedisService) Delete(key string) error {
	conn := this.GetRedis()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

func (this *RedisService) IsExist(key string) bool {
	conn := this.GetRedis()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (this *RedisService) HMGet(key string, fields ...string) (reply interface{}, err error) {
	redisConn := this.GetRedis()
	defer redisConn.Close()
	args := make([]string, 0)
	args = append(args, key)
	for _, v := range fields {
		args = append(args, v)
	}
	return redisConn.Do("hmget", args)
}

func (this *RedisService) Expire(key string, timeout time.Duration) (reply interface{}, err error) {
	conn := this.GetRedis()
	defer conn.Close()
	fmt.Println(key, timeout.Seconds())
	seconds := int(timeout.Seconds())
	reply, err = conn.Do("expire", key, seconds)
	return
}

func (this *RedisService) IncrWithExpire(key string, expire int) (count int64, err error) {
	conn := this.GetRedis()
	defer conn.Close()
	scriptString := `local current
current = redis.call("incr",KEYS[1])
if tonumber(current) == 1 then
    redis.call("expire",KEYS[1], ARGV[1])
end
return current`
	script := redis.NewScript(1, scriptString)
	reply, err := script.Do(conn, key, expire)
	count = utils.ToInt64(reply)
	return
}

func (this *RedisService) Decr(key string) (value int64, err error) {
	redisConn := this.GetRedis()
	defer redisConn.Close()
	reply, err := redisConn.Do("decr", key)
	value = cache.GetInt64(reply)
	return
}
