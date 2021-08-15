package business

import (
	"dongchamao/global/utils"
	"dongchamao/services"
	"net/http"
	"regexp"
	"strings"
)

type DouyinBusiness struct {
}

func NewDouyinBusiness() *DouyinBusiness {
	return new(DouyinBusiness)
}

//短url还原解析
func (this *DouyinBusiness) ParseDyShortUrl(url string) string {
	filterUrl := utils.ParseDyVideoShare(url)
	if filterUrl != "" && filterUrl != url {
		url = filterUrl
	}
	url = strings.TrimSpace(url)
	//判断是否短网址,之后加入缓存
	pattern := `^(http|https):\/\/v\.douyin\.com\/.*?`
	reg := regexp.MustCompile(pattern)
	returl := ""
	if reg.MatchString(url) == true {
		redisService := services.NewRedisService()
		returl = redisService.Hget("douyin:shorturl:hashmap", url)
		if returl == "" {
			client := &http.Client{}
			request, _ := http.NewRequest("GET", url, nil)
			response, err := client.Do(request)
			if err != nil {
				return ""
			}
			defer response.Body.Close()
			returl = response.Request.URL.String()
			redisService.Hset("douyin:shorturl:hashmap", url, returl)
		}
		return returl
	} else {
		return url
	}
}
