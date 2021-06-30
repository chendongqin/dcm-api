package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

type RClient struct {
	db         string
	clientPool *redis.Pool
}

func (rc *RClient) getPool() *redis.Pool {
	if rc.clientPool != nil {
		return rc.clientPool
	} else {
		return nil
	}
}

func (rc *RClient) SetPool(db string, p *redis.Pool) {
	rc.clientPool = p
	rc.db = db
}

// 这边会有个redis.ErrNil,是string里抛出的，不是redis问题
func (rc *RClient) GetInt(key string) (int, error) {

	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("GET", key)
		val, err := redis.Int(v, err)
		return val, err
	} else {
		return 0, client.Err()
	}

}

func (rc *RClient) HGetInt(key string, subKey string) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("HGET", key, subKey)
		val, _ := redis.Int(v, err)
		return val, err
	} else {
		return 0, client.Err()
	}

}

func (rc *RClient) MGetInt(keys []string) ([]int, error) {
	if len(keys) == 0 {
		return []int{}, nil
	}
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		var args []interface{}
		for _, v := range keys {
			args = append(args, v)
		}
		v, err := client.Do("MGET", args...)
		val, _ := redis.Ints(v, err)
		return val, err
	} else {
		return []int{}, client.Err()
	}
}

func (rc *RClient) Incr(key string, expire int) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("INCR", key)
		val, _ := redis.Int(v, err)
		if val <= 1 && expire != -1 {
			client.Do("EXPIRE", key, expire)
		}
		return val, nil
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) Decr(key string) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("DECR", key)
		val, _ := redis.Int(v, err)
		return val, nil
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) IncrKeys(keys []string, expire int) error {

	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		for _, key := range keys {
			v, err := client.Do("INCR", key)
			val, _ := redis.Int(v, err)
			if val <= 1 {
				client.Do("EXPIRE", key, expire)
			}
		}
		return nil
	} else {
		return client.Err()
	}
}

// 这边会有个redis.ErrNil,是string里抛出的，不是redis问题
func (rc *RClient) GetString(key string) (string, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("GET", key)
		val, err := redis.String(v, err)
		return val, err
	} else {
		return "", client.Err()
	}
}

func (rc *RClient) MGetString(keys []string) ([]string, error) {
	if len(keys) == 0 {
		return []string{}, nil
	}
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		var args []interface{}
		for _, v := range keys {
			args = append(args, v)
		}
		v, err := client.Do("MGET", args...)

		val, _ := redis.Strings(v, err)
		return val, err
	} else {
		return []string{}, client.Err()
	}
}

func (rc *RClient) Set(params ...interface{}) (bool, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("SET", params...)

		if v == "OK" && err == nil {
			return true, nil
		}
		return false, err
	} else {
		return false, client.Err()
	}
}

func (rc *RClient) SetEx(params ...interface{}) (bool, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("SETEX", params...)
		if v == "OK" && err == nil {
			return true, nil
		}
		return false, err
	} else {
		return false, client.Err()
	}
}

func (rc *RClient) SetString(key, data string, exp int) (bool, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("SET", key, data, "EX", exp)
		if v == "OK" && err == nil {
			return true, nil
		}
		return false, err
	} else {
		return false, client.Err()
	}
}

func (rc *RClient) SetInt(key string, data int, exp int) (bool, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("SET", key, data, "EX", exp)
		if v == "OK" && err == nil {
			return true, nil
		}
		return false, err
	} else {
		return false, client.Err()
	}
}

func (rc *RClient) Hincrby(key, key2 string, step, expire int) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("HINCRBY", key, key2, step)
		val, _ := redis.Int(v, err)
		if val <= 1 {
			client.Do("EXPIRE", key, expire)
		}
		return val, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) Hset(key string, field string, value string) (bool, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("HSET", key, field, value)
		if v == "OK" && err == nil {
			return true, nil
		}
		return false, err
	} else {
		return false, client.Err()
	}
}

func (rc *RClient) Hmset(key string, dataMap map[string]string, expire int) (bool, error) {
	if len(dataMap) == 0 {
		return false, nil
	}
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		args := redis.Args{}.Add(key).AddFlat(dataMap)
		v, err := client.Do("HMSET", args...)
		if v == "OK" && err == nil {
			client.Do("EXPIRE", key, expire)
			return true, err
		}
		return false, err
	} else {
		return false, client.Err()
	}
}

func (rc *RClient) Hmget(key string, values []string) (map[string]string, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		args := redis.Args{}.Add(key).AddFlat(values)
		res, err := redis.Strings(client.Do("HMGET", args...))

		stringMap1 := map[string]string{}
		for i, hashv := range res {
			if hashv != "" {
				valKey := values[i]
				stringMap1[valKey] = hashv
			}
		}

		return stringMap1, err
	} else {
		return map[string]string{}, client.Err()
	}
}

func (rc *RClient) HGetString(key, value string) (string, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		args := redis.Args{}.Add(key).Add(value)
		v, err := client.Do("HGET", args...)
		val, err := redis.String(v, err)
		return val, err
	} else {
		return "", client.Err()
	}
}

func (rc *RClient) Keys(keyPattern string) ([]string, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		args := redis.Args{}.Add(keyPattern)
		res, err := redis.Strings(client.Do("KEYS", args...))
		return res, err
	} else {
		return nil, client.Err()
	}
}

func (rc *RClient) Hmgetall(key string) (map[string]string, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.StringMap(client.Do("HMGETALL", redis.Args{}.Add(key)...))

		return v, err
	} else {
		return map[string]string{}, client.Err()
	}
}

func (rc *RClient) LPush(key string, values interface{}) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("LPUSH", redis.Args{}.Add(key).AddFlat(values)...))
		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) RPopInt(key string) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("RPOP", key))
		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) RPopBytes(key string) ([]byte, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Bytes(client.Do("RPOP", key))
		return v, err
	} else {
		return []byte{}, client.Err()
	}
}

func (rc *RClient) LLen(key string) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("LLEN", key))
		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) SIsMember(key string, value interface{}) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("SISMEMBER", key, value))

		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) SMembers(keys string) ([]string, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("SMEMBERS", keys)
		val, e := redis.Values(v, err)
		value := make([]string, 0)
		if e == nil {
			for _, v := range val {
				value = append(value, fmt.Sprintf("%s", v.([]byte)))
			}
		}
		return value, e
	} else {
		return []string{}, client.Err()
	}
}

func (rc *RClient) Exists(key string) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("EXISTS", key))

		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) Expire(key string, ttl int) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("EXPIRE", key, ttl))

		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) TTL(key string) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("TTL", key))

		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) SAdd(key string, values ...interface{}) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("SADD", redis.Args{}.Add(key).Add(values...)...))
		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) SCard(key string) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("SCARD", key))
		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) SRem(key string, values interface{}) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("SREM", key, values))
		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) Del(key string) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Int(client.Do("DEL", key))
		return v, err
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) IncrBy(key string, add, expire int) (int, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("INCRBY", key, add)
		val, _ := redis.Int(v, err)
		if val == add && expire != -1 {
			client.Do("EXPIRE", key, expire)
		}
		return val, nil
	} else {
		return 0, client.Err()
	}
}

func (rc *RClient) ZAdd(key string, values map[float64]interface{}, expire int) (bool, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := client.Do("ZADD", redis.Args{}.Add(key).AddFlat(values)...)
		val, _ := redis.Int(v, err)
		if val == 0 && expire != -1 {
			client.Do("EXPIRE", key, expire)
		}
		return false, err
	} else {
		return false, client.Err()
	}
}

func (rc *RClient) ZRevRange(key string, start, end int) ([]string, error) {
	client := rc.clientPool.Get()
	if client.Err() == nil {
		defer client.Close()
		v, err := redis.Strings(client.Do("ZREVRANGE", key, start, end))
		return v, err
	} else {
		return nil, client.Err()
	}
}

var updateRedisLock sync.RWMutex
var redisServerList = make(map[string]*redis.Pool)

type Options struct {
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
	Host        string
	PassWd      string
	DbName      string
	DB          int
}

func Redis(dbName string) *RClient {
	var rc = &RClient{}
	if pool, ok := redisServerList[dbName]; ok {
		rc.SetPool(dbName, pool)
	} else {
		return nil
	}
	return rc
}

func InitRedis(dbName string, opt Options) {
	updateRedisLock.Lock()
	redisServerList[dbName] = func(opt Options) *redis.Pool {
		return &redis.Pool{
			MaxIdle:     opt.MaxIdle,
			MaxActive:   opt.MaxActive,
			Wait:        true,
			IdleTimeout: time.Duration(opt.IdleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", opt.Host,
					redis.DialPassword(opt.PassWd),
					redis.DialDatabase(opt.DB))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
	}(opt)
	updateRedisLock.Unlock()
}
