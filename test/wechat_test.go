package test

import (
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"testing"
)

func TestWxAccessToken(t *testing.T) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory() //TODO 建议改成 REDIS
	cfg := &config.Config{
		AppID:          "wx6392924bcb25fd24",
		AppSecret:      "a0bc0bcdbb77219655ff8aff1b110559",
		Token:          "7B6DB2CF18CC177BBFB5ED72B3644899",
		EncodingAESKey: "Qgnn0UH9oz5rNHZIex27UkY4Xxu3daMUFbHvLaaaRwe",
		Cache:          memory,
	}
	WxOfficial := wc.GetOfficialAccount(cfg)
	accToken, _ := WxOfficial.GetAccessToken()
	fmt.Println(accToken)
}
