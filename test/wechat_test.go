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
	code := "001ZTQFa1DSvFB0rWwFa17BQOl3ZTQFI"
	userInfo, _ := business.NewWxAppBusiness().AppLogin(code)
	fmt.Println(userInfo)
	//fmt.Println(err)
}

//获取微信公众号
func TestWxMenu(t *testing.T){
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
	menuInfo := WxOfficial.GetMenu()
	fmt.Println(menuInfo.GetMenu())
}

//获取微信公众号
func TestDelWxMenu(t *testing.T){
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
	menuInfo := WxOfficial.GetMenu()
	fmt.Println(menuInfo.DeleteMenu())
}

func TestSetWxMenu(t *testing.T){
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
	menuInfo := WxOfficial.GetMenu()
	fmt.Println(menuInfo.SetMenuByJSON(`{"button":[{"type":"click","name":"关于我们","key":" ABOUT_UD"},{"name":"主菜单","sub_button":[{"type":"view","name":"搜索","url":"http://www.soso.com/"},{"type":"click","name":"赞一下我们","key":"V1001_GOOD"}]}]}`))
}
