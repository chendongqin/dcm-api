package redis

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	redisOpt := Options{
		MaxIdle:   5,
		MaxActive: 5,
		Host:      "127.0.0.1:6379",
		PassWd:    "",
		DB:        1,
	}
	InitRedis("default", redisOpt)

	c := Redis("default")
	// fmt.Println(c.LPush("lpushtest", []string{"zzzza", "zzzzb "}))
	// fmt.Println(count, err)
	// fmt.Println(c.Set("hhhhhasdfasd", "asdfasdfasd", "NX", "EX", 60))

	// fmt.Println(c.SIsMember("fqb:device:install:894394", 1197))
	//c.SetString("aa", "asdfasdfads", 60)
	//c.Del("aa")
	//fmt.Println(c.GetString("aa"))

	/*
			intS := make([]int, 0)
			for i := 1; i < 233; i++ {
				intS = append(intS, i)
			}

			j := 0
			for j < len(intS) {
				if j+100 < len(intS) {
					fmt.Println(j, j+100)
					fmt.Println(c.LPush("lpushtest1z", intS[j:j+100]))
					j = j + 100
				} else {
					fmt.Println(j)
					fmt.Println(c.LPush("lpushtest1z", intS[j:]))
					j = len(intS)
				}
			}


		fmt.Println(c.SAdd("saddaaa", 1))
		fmt.Println(c.SAdd("saddaaa", 1))
	*/
	fmt.Println(c.Decr("test_decr"))
	fmt.Println(c.Exists("test_decr111"))
}
