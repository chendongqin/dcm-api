package utils

import (
	"dongchamao/global"
	"dongchamao/services/mutex"
	//"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	json "github.com/json-iterator/go"
	"math"
	"reflect"
	"time"
)

// Example:
// key := "test"
// cacheHelper := utils.NewCacheHelper(key)
// obj := make([]string) // obj can be any type struct or slice.
// source := func () interface{} { ... return data or err }
// cacheHelper.SetSource(source).Get(&obj)

// Key: 缓存Key
// Expire: 缓存过期时间
// ReadTimeout: 读取缓存超时时间
// LockTimeout: 锁超时时间
// Source: 获取源数据方法

type CacheHelper struct {
	Key         string
	Expire      time.Duration
	ReadTimeout time.Duration
	LockTimeout time.Duration
	Source      func(c *CacheHelper) interface{}
}

func (helper *CacheHelper) SetKey(key string) *CacheHelper {
	helper.Key = key
	return helper
}

func (helper *CacheHelper) SetSource(source interface{}) *CacheHelper {
	if s, ok := source.(func() interface{}); ok {
		helper.Source = func(c *CacheHelper) interface{} {
			return s()
		}
	} else if s, ok := source.(func(c *CacheHelper) interface{}); ok {
		helper.Source = s
	} else {
		helper.Source = func(c *CacheHelper) interface{} {
			logs.Error("[cache] source nothing to do")
			return nil
		}
	}
	return helper
}

func (helper *CacheHelper) SetExpire(expire time.Duration) *CacheHelper {
	helper.Expire = expire / time.Second
	return helper
}

func (helper *CacheHelper) SetLockTimeout(lockTimeout time.Duration) *CacheHelper {
	helper.LockTimeout = lockTimeout / time.Second
	return helper
}

func (helper *CacheHelper) SetReadTimeout(readTimeout time.Duration) *CacheHelper {
	helper.ReadTimeout = readTimeout / time.Second
	return helper
}

func (helper *CacheHelper) GetCacheString() string {
	return global.Cache.Get(helper.Key)
}

func (helper *CacheHelper) Get(result interface{}, force ...bool) (err error) {
	if helper.Source == nil {
		err = errors.New("the data source callback is nil")
		logs.Error(err)
		return
	}
	data := helper.GetCacheString()
	if data != "" {
		logs.Debug("get data from the cache successfully.[%s]", helper.Key)
		err = helper.unmarshalCache(data, result)
	} else {
		if len(force) > 0 && force[0] {
			err = helper.getSource(result)
		} else {
			err = helper.acquireSource(result)
		}
	}
	return
}

func (helper *CacheHelper) acquireSource(result interface{}) (err error) {
	lockToken := GetRandomString(16)
	var lock *mutex.Lock
	var ok bool
	lock, ok, err = mutex.TryLockWithTimeout(global.Cache.GetInstance().(redis.Conn), helper.Key+":lock", lockToken, int(helper.LockTimeout))
	if lock != nil && ok {
		defer lock.Unlock()
	}
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Debug("acquire the lock successfully, try to get data from source callback")
	// acquire the lock successfully, try to get data from source callback
	if ok {
		err = helper.getSource(result)
		return
	}
	logs.Debug("failure to get the lock, and retrying to get data from the cache.")
	// failure to get the lock, and retrying to get data from the cache.
	maxRetryTimes := int(math.Floor(float64(helper.ReadTimeout*1000) / float64(100)))
	logs.Debug("max retry times", maxRetryTimes)
	for i := 0; i < maxRetryTimes; i++ {
		time.Sleep(time.Millisecond * 100)
		data := helper.GetCacheString()
		if data == "" {
			continue
		}
		err = helper.unmarshalCache(data, result)
		return
	}
	return &global.TimeoutError{}
}

func (helper *CacheHelper) unmarshalCache(data string, result interface{}) (err error) {
	err = json.Unmarshal([]byte(data), result)
	return
}

func (helper *CacheHelper) getSource(result interface{}) (err error) {
	logs.Debug("get data from source callback")
	temp := helper.Source(helper)
	if sourceErr, ok := temp.(error); ok {
		logs.Error("[cache] the source return an error, err: %s", sourceErr)
		err = sourceErr
		return
	}
	if temp == nil {
		return
	}
	var jsonByte []byte
	jsonByte, err = json.Marshal(temp)
	err = json.Unmarshal(jsonByte, result)
	dataSize := 1
	// 如果是slice, 则获取实际长度
	if reflect.TypeOf(temp).Kind() == reflect.Slice {
		dataSize = helper.Size(temp)
	}
	// slice长度大于0时才写入缓存
	expire := time.Duration(86400)
	if helper.Expire > 0 {
		expire = helper.Expire
	}
	if dataSize <= 0 {
		expire = time.Duration(60)
	}
	err = global.Cache.SetMap(helper.Key, temp, expire)
	return
}

func (helper CacheHelper) Size(a interface{}) int {
	ins := reflect.ValueOf(a)
	return ins.Len()
}

func NewCacheHelper(key string) *CacheHelper {
	return &CacheHelper{
		Key:         key,
		Expire:      time.Duration(86400),
		ReadTimeout: time.Duration(3),
		LockTimeout: time.Duration(3),
		Source:      nil,
	}
}

/*
	缓存不存在时阻塞等待，等待超时返回错误
*/
func SimpleCacheGet(key string, source interface{}, result interface{}, expire ...int) (err error) {
	cacheHelper := NewCacheHelperQuick(key, expire...)
	err = cacheHelper.SetSource(source).
		Get(result)
	return
}

/**
缓存不存在时不等待，直接从source中获取数据
*/
func SimpleCacheGetForce(key string, source interface{}, result interface{}, expire ...int) (err error) {
	cacheHelper := NewCacheHelperQuick(key, expire...)
	err = cacheHelper.SetSource(source).
		Get(result, true)
	return
}

func NewCacheHelperQuick(key string, expire ...int) *CacheHelper {
	cacheHelper := NewCacheHelper(key)
	if len(expire) > 0 {
		cacheHelper.SetExpire(time.Duration(expire[0]) * time.Second)
		if len(expire) > 1 {
			cacheHelper.SetReadTimeout(time.Duration(expire[1]) * time.Second)
			if len(expire) > 2 {
				cacheHelper.SetLockTimeout(time.Duration(expire[2]) * time.Second)
			}
		}
	}
	return cacheHelper
}
