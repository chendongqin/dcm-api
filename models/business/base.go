package business

import (
	"dongchamao/structinit/repost/dy"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

const (
	VoUserUniqueTokenPlatformWeb           = 1
	VoUserUniqueTokenPlatformH5            = 2
	VoUserUniqueTokenPlatformWxMiniProgram = 3
	VoUserUniqueTokenPlatformApp           = 4
	VoUserUniqueTokenPlatformWap           = 5
)

type AuthorCate struct {
	First  map[string]string   `json:"first"`
	Second []map[string]string `json:"second"`
}

func GetAppPlatFormIdWithAppId(appId int) int {
	switch appId {
	case 10000:
		return VoUserUniqueTokenPlatformWeb
	case 10001:
		return VoUserUniqueTokenPlatformH5
	case 10002:
		return VoUserUniqueTokenPlatformWxMiniProgram
	case 10003, 10004:
		return VoUserUniqueTokenPlatformApp
	case 10005:
		return VoUserUniqueTokenPlatformWap
	}
	return 0
}

func GetAge(startTime string) int {
	var Age int64
	var pslTime string
	if strings.Index(startTime, ".") != -1 {
		pslTime = "2006.01.02"
	} else if strings.Index(startTime, "-") != -1 {
		pslTime = "2006-01-02"
	} else {
		pslTime = "2006/01/02"
	}
	t1, err := time.ParseInLocation(pslTime, startTime, time.Local)
	t2 := time.Now()
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		Age = diff / (3600 * 365 * 24)
		return int(Age)
	} else {
		return int(Age)
	}
}

func DealChartInt(charts []int, chartNum int) []int {
	nums := len(charts)
	if nums <= chartNum {
		return charts
	}
	newCharts := make([]int, 0)
	randNum := nums / chartNum
	if randNum < 1 {
		randNum = 1
	}
	var begin int
	for i := 0; i < chartNum; i++ {
		begin = i * randNum
		if begin >= nums {
			break
		}
		newCharts = append(newCharts, charts[begin])
	}
	return newCharts
}

func DealChartInt64(charts []int64, chartNum int) []int64 {
	nums := len(charts)
	if nums <= chartNum {
		return charts
	}
	newCharts := make([]int64, 0)
	randNum := nums / chartNum
	if randNum < 1 {
		randNum = 1
	}
	var begin int
	for i := 0; i < chartNum; i++ {
		begin = i * randNum
		if begin >= nums {
			break
		}
		newCharts = append(newCharts, charts[begin])
	}
	return newCharts
}

func DealChartFloat64(charts []float64, chartNum int) []float64 {
	nums := len(charts)
	if nums <= chartNum {
		return charts
	}
	newCharts := make([]float64, 0)
	randNum := nums / chartNum
	if randNum < 1 {
		randNum = 1
	}
	var begin int
	for i := 0; i < chartNum; i++ {
		begin = i * randNum
		if begin >= nums {
			break
		}
		newCharts = append(newCharts, charts[begin])
	}
	return newCharts
}

func DealAuthorCateJson(authorCateJson string) []dy.DyCate {
	authorCateMap := map[string][]AuthorCate{}
	_ = jsoniter.Unmarshal([]byte(authorCateJson), &authorCateMap)
	cates := make([]AuthorCate, 0)
	if v, ok := authorCateMap["tag"]; ok {
		cates = v
	}
	dyAuthorCate := make([]dy.DyCate, 0)
	for _, v := range cates {
		firName := ""
		for _, name := range v.First {
			firName = name
		}
		item := dy.DyCate{
			Name:    firName,
			SonCate: []dy.DyCate{},
		}
		for _, s := range v.Second {
			for _, name := range s {
				item.SonCate = append(item.SonCate, dy.DyCate{
					Name:    name,
					SonCate: []dy.DyCate{},
				})
			}
		}
		dyAuthorCate = append(dyAuthorCate, item)
	}
	return dyAuthorCate
}
