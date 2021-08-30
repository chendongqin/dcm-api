package test

import (
	"dongchamao/business"
	"dongchamao/global/consistent"
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"strings"
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

//获取用户信息 (公众号)
func TestUserInfo(t *testing.T) {
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
	userInfo, _ := WxOfficial.GetUser().GetUserInfo("oYywQ54At-8F_hTmOZMa40cG9ygA")
	fmt.Println(userInfo)
}

func TestUserStr(t *testing.T) {
	eventKey := "qrscene_qrlogin:ca46379c043766ae95e19be66105e879"
	fmt.Println(strings.Contains(eventKey, consistent.WECHAT_QR_LOGIN))
}

//客户端微信登录
func TestWxApp(t *testing.T) {
	code := "081wfm0w3lxNYW2BgJ1w3CaWoR0wfm04"
	userInfo, err := business.NewWxAppBusiness().AppLogin(code)
	fmt.Println(userInfo)
	fmt.Println(err)
}
