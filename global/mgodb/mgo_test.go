package mgodb

import (
	"fmt"
	"testing"
	"time"
)

type RepeatAppIdfa struct {
	Idfa    string `bson:"idfa"`
	Channel string `bson:"channel"`
	AppId   int    `bson:"appId"`
}

func TestFind(t *testing.T) {
	InitMongo(Options{
		Addr:     "127.0.0.1:27017",
		Database: "chanshike",
		Username: "",
		Password: "",
		Idle:     10,
	})
	condition_1 := map[string]interface{}{
		"idfa": map[string]interface{}{
			"$in": []string{"1956ABEE-B112-4FC7-8696-207C6F4E3F90"},
		},
		"appId": 1001153553,
	}
	reSlice := make([]RepeatAppIdfa, 0)
	err := NewMgoSelector().Table("repeat_app_idfa").Where(condition_1).Select(&reSlice)
	fmt.Println(err, reSlice)
}

type ActiveLog struct {
	Idfa    string    `bson:"idfa"`
	Udid    string    `bson:"udid"`
	Channel string    `bson:"channel"`
	AppId   int       `bson:"appId"`
	Created time.Time `bson:"created"`
}

func TestUpinsert(t *testing.T) {
	InitMongo(Options{
		Addr:     "127.0.0.1:27017",
		Database: "chanshike",
		Username: "",
		Password: "",
		Idle:     10,
	})
	condition_1 := map[string]interface{}{
		"idfa":  "1956ABEE-B112-4FC7-8696-207C6F4E3F90",
		"appId": 1461695441,
	}
	reSlice := make([]ActiveLog, 0)
	err := NewMgoSelector().Table("active_log").Where(condition_1).Select(&reSlice)
	fmt.Println(err, reSlice)
	fmt.Println(NewMgoSelector().Table("active_log").UpInsert(condition_1, ActiveLog{
		Idfa:    "1956ABEE-B112-4FC7-8696-207C6F4E3F90",
		AppId:   1461695441,
		Created: time.Now(),
	}))
}

func TestRemove(t *testing.T) {
	InitMongo(Options{
		Addr:     "127.0.0.1:27017",
		Database: "chanshike",
		Username: "",
		Password: "",
		Idle:     10,
	})
	condition_1 := map[string]interface{}{
		"created": map[string]interface{}{
			"$lte": time.Now().AddDate(0, 0, -2),
		},
	}
	err := NewMgoSelector().Table("active_log").RemoveOne(condition_1)
	fmt.Println(err)
}

type RepeatIdfaSet struct {
	Idfa  string `bson:"idfa"`
	AppId []int  `bson:"appId"`
}

func TestSet(t *testing.T) {
	InitMongo(Options{
		Addr:     "127.0.0.1:27017",
		Database: "chanshike",
		Username: "",
		Password: "",
		Idle:     10,
	})
	condition1 := map[string]interface{}{
		"idfa": "540784D3-D745-4E17-A69C-0567BEB93DFA",
		"appId": map[string]interface{}{
			"$in": []int{539988421},
		},
	}
	reSlice := make([]RepeatIdfaSet, 0)
	err := NewMgoSelector().Table("repeat_idfa_set").Where(condition1).Select(&reSlice)
	fmt.Println(reSlice, err)
}

func TestAddToSet(t *testing.T) {
	InitMongo(Options{
		Addr:     "127.0.0.1:27017",
		Database: "chanshike",
		Username: "",
		Password: "",
		Idle:     10,
	})
	condition_1 := map[string]interface{}{
		"idfa": "1956ABEE-B112-4FC7-8696-207C6F4E3F90",
	}
	reSlice := make([]ActiveLog, 0)
	err := NewMgoSelector().Table("repeat_idfa_set").Where(condition_1).Select(&reSlice)
	fmt.Println(err, reSlice)
	fmt.Println(NewMgoSelector().Table("repeat_idfa_set").UpInsert(condition_1, map[string]interface{}{
		"$addToSet": map[string]int{
			"appId": 539988422,
		},
	}))
}
