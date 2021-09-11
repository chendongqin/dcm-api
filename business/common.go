package business

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/hbase"
	"dongchamao/models/dcm"
	"dongchamao/models/repost/dy"
	"dongchamao/services"
	"encoding/base64"
	"errors"
	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"regexp"
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

const (
	EsMaxShowNum          = 10000
	DyJewelBaseShowNum    = 5000
	DyJewelBaseMinShowNum = 10
	DyJewelRankShowNum    = 1500
	DyRankMinShowNum      = 5
)

const (
	AwemeUrl     = "https://www.douyin.com/video/"
	LiveShareUrl = "https://www.iesdouyin.com/share/live/"
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

func DealAuthorLiveTags() {
	data, _ := hbase.GetAuthorLiveTags()
	tagsMap := map[string]string{}
	for _, v := range data {
		tags := strings.Split(v.Tags, "_")
		for _, tag := range tags {
			tagsMap[tag] = tag
		}
	}
	for _, tag := range tagsMap {
		tagM := dcm.DyAuthorLiveTags{}
		exist, _ := dcm.GetBy("name", tag, &tagM)
		if exist {
			continue
		}
		tagM.Name = tag
		dcm.Insert(nil, &tagM)
	}
	return
}

func GetConfig(keyName string) string {
	cacheKey := cache.GetCacheKey(cache.ConfigKeyCache, keyName)
	cacheData := global.Cache.Get(cacheKey)
	if cacheData == "" {
		var configJson dcm.DcConfigJson
		exist, err := dcm.GetBy("key_name", keyName, &configJson)
		if exist && err == nil {
			cacheData = configJson.Value
		}
		_ = global.Cache.Set(cacheKey, cacheData, 1800)
	}
	return cacheData
}

func UserActionLock(active string, userData string, lockTime time.Duration) bool {
	memberKey := cache.GetCacheKey(cache.UserActionLock, active, userData)
	if global.Cache.Get(memberKey) != "" {
		return false
	}
	_ = global.Cache.Set(memberKey, "1", lockTime)
	return true
}

//短url还原解析
func ParseDyShortUrl(url string) (string, bool) {
	url = strings.TrimSpace(url)
	//判断是否短网址,之后加入缓存
	pattern := `^(http|https):\/\/v\.douyin\.com\/.*?`
	reg := regexp.MustCompile(pattern)
	returl := ""
	if reg.MatchString(url) == true {
		redisService := services.NewRedisService()
		returl = redisService.Hget("douyin:shorturl:hashmap", url)
		if returl == "" {
			client := &http.Client{CheckRedirect: nil}
			request, _ := http.NewRequest("GET", url, nil)
			response, err := client.Do(request)
			if err != nil {
				return "", false
			}
			defer response.Body.Close()
			returl = response.Request.Response.Request.URL.Path
			if len(returl) == 0 {
				return "", false
			}
			redisService.Hset("douyin:shorturl:hashmap", url, returl)
		}
		return returl, true
	} else {
		logs.Info("[短链转换失败][%s]", url)
		return url, false
	}
}

//id加密
func IdEncrypt(id string) string {
	if global.IsDev() {
		return id
	}
	if id == "" || id == "0" {
		return ""
	}
	key := []byte("dwVRjLVUN4RMGAKSEvuvPV696PKrEuRT")
	idByte := []byte(id)
	str, err := utils.AesEncrypt(idByte, key)
	if err != nil {
		return ""
	}
	//restful路由避免错误
	return "==" + strings.ReplaceAll(base64.StdEncoding.EncodeToString(str), "/", "*")
}

//id解密
func IdDecrypt(id string) string {
	if id == "" {
		return ""
	}
	if strings.Index(id, "==") != 0 {
		if global.IsDev() {
			return id
		}
		return ""
	}
	id = strings.Replace(id, "==", "", 1)
	key := []byte("dwVRjLVUN4RMGAKSEvuvPV696PKrEuRT")
	//restful路由避免错误
	id = strings.ReplaceAll(id, "*", "/")
	s, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return ""
	}
	str, err := utils.AesDecrypt(s, key)
	if err != nil {
		return ""
	}
	return string(str)
}

//json加密
func JsonEncrypt(jsonData interface{}) string {
	jsonByte := []byte{}
	switch result := jsonData.(type) {
	case string:
		jsonByte = []byte(result)
	default:
		jsonByte, _ = jsoniter.Marshal(jsonData)
	}
	if len(jsonByte) == 0 {
		return ""
	}
	key := []byte("LFROPI0K0w5JVauUBLEexvvDTHxaxCZL")
	str, err := utils.AesEncrypt(jsonByte, key)
	if err != nil {
		return ""
	}
	//restful路由避免错误
	return base64.StdEncoding.EncodeToString(str)
}

//json解密
func JsonDecrypt(json string) string {
	if json == "" {
		return ""
	}
	key := []byte("LFROPI0K0w5JVauUBLEexvvDTHxaxCZL")
	//restful路由避免错误
	s, err := base64.StdEncoding.DecodeString(json)
	if err != nil {
		return ""
	}
	str, err := utils.AesDecrypt(s, key)
	if err != nil {
		return ""
	}
	return string(str)
}

func WitheUsername(username string) error {
	cacheKey := cache.GetCacheKey(cache.ConfigKeyCache, "begin_open_phone")
	cacheData := global.Cache.Get(cacheKey)
	list := make([]string, 0)
	if cacheData != "" {
		_ = jsoniter.Unmarshal([]byte(cacheData), &list)
	} else {
		var configJson dcm.DcConfigJson
		_, _ = dcm.GetBy("key_name", "begin_open_phone", &configJson)
		list = strings.Split(configJson.Value, ",")
		listByte, _ := jsoniter.Marshal(list)
		global.Cache.Set(cacheKey, string(listByte), 300)
	}
	newList := make([]string, 0)
	for _, v := range list {
		if v != "" {
			newList = append(newList, v)
		}
	}
	if len(newList) == 0 {
		return nil
	}
	if !utils.InArrayString(username, list) {
		return errors.New("您没有权限~")
	}
	return nil
}

//处理脏数据
func DealIncDirtyInt64Chart(chart []int64) []int64 {
	chart = utils.ReverseInt64Arr(chart)
	lenNum := len(chart)
	if lenNum == 2 {
		if chart[0] < chart[1] {
			chart[1] = chart[0]
		}
	} else if lenNum > 2 {
		for k, v := range chart {
			if k == 0 {
				continue
			} else if k == lenNum-1 {
				if v > chart[k-1] {
					chart[k] = chart[k-1]
				}
				break
			}
			if chart[k+1] > v && chart[k-1] > chart[k+1] {
				chart[k] = chart[k-1]
				continue
			}
		}
	}
	chart = utils.ReverseInt64Arr(chart)
	return chart
}

func DealIncDirtyFloat64Chart(chart []float64) []float64 {
	chart = utils.ReverseFloat64Arr(chart)
	lenNum := len(chart)
	if lenNum == 2 {
		if chart[0] < chart[1] {
			chart[1] = chart[0]
		}
	} else if lenNum > 2 {
		for k, v := range chart {
			if k == 0 {
				continue
			} else if k == lenNum-1 {
				if v > chart[k-1] {
					chart[k] = chart[k-1]
				}
				break
			}
			if chart[k+1] > v && chart[k-1] > chart[k+1] {
				chart[k] = chart[k-1]
				continue
			}
		}
	}
	chart = utils.ReverseFloat64Arr(chart)
	return chart
}
