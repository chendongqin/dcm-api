package cache

import "fmt"

func GetCacheKey(key KeyName, a ...interface{}) string {
	return fmt.Sprintf(string(key), a...)
}
