package mutex

import (
	"dongchamao/global"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

type Lock struct {
	resource string
	token    string
	conn     redis.Conn
	timeout  int
}

func (lock *Lock) tryLock() (ok bool, err error) {
	startTime := time.Now()
	defer func() {
		logs.Debug("[redis-lock] [tryLock] [resource:%s] [took:%s]", lock.key(), time.Since(startTime))
	}()
	_, err = redis.String(lock.do("SET", lock.key(), lock.token, "EX", int(lock.timeout), "NX"))
	if err == redis.ErrNil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (lock *Lock) Unlock() (err error) {
	startTime := time.Now()
	defer func() {
		logs.Debug("[redis-lock] [UnLock] [resource:%s] [took:%s]", lock.key(), time.Since(startTime))
	}()
	var str string
	str, err = redis.String(lock.do("get", lock.key()))
	if str == lock.token {
		_, err = lock.do("del", lock.key())
	} else {
		err = errors.New("unlock failed")
	}
	return
}

func (lock *Lock) key() string {
	return fmt.Sprintf("redislock:%s", lock.resource)
}

func (lock *Lock) AddTimeout(exTime int64) (ok bool, err error) {
	ttlTime, err := redis.Int64(lock.do("TTL", lock.key()))
	if err != nil {
		return
	}
	if ttlTime > 0 {
		_, err := redis.String(lock.do("SET", lock.key(), lock.token, "EX", int(ttlTime+exTime)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func TryLock(conn redis.Conn, resource string, token string, DefaulTimeout int) (lock *Lock, ok bool, err error) {
	return TryLockWithTimeout(conn, resource, token, DefaulTimeout)
}

func TryLockWithTimeout(conn redis.Conn, resource string, token string, timeout int) (lock *Lock, ok bool, err error) {
	lock = &Lock{resource, token, conn, timeout}

	ok, err = lock.tryLock()

	if err != nil {
		lock = nil
	}

	return
}

func (lock *Lock) WaitForUnlock(timeout time.Duration) (bool, error) {
	times := 0
	delay := time.Millisecond * 50
	maxTimes := int(timeout / delay)
	for {
		str, err := redis.String(lock.do("get", lock.key()))
		if err != nil {
			return false, err
		}
		if str == "" {
			return true, nil
		}
		times++
		if times > maxTimes {
			return false, errors.New("wait timeout")
		}
		time.Sleep(delay)
	}
}

func (lock *Lock) Retry() (bool, error) {
	return lock.tryLock()
}

func (lock *Lock) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = lock.associate(args[0])
	c := global.Cache.GetInstance().(redis.Conn)
	defer c.Close()
	return c.Do(commandName, args...)
}

// associate with config key.
func (lock *Lock) associate(originKey interface{}) string {
	return fmt.Sprintf("%s", originKey)
}
