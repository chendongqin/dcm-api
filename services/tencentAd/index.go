package tencent_ad

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"github.com/silenceper/wechat/v2/util"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"time"
)

type AccessToken struct {
	TAds           *ads.SDKClient
	AccessToken    string
	ClientId       int64
	ClientSecret   string
	GrantType      string
	OauthTokenOpts *api.OauthTokenOpts
}

func (e *AccessToken) Init() {
	e.TAds = ads.Init(&config.SDKConfig{
		IsDebug: true,
	})
	e.TAds.UseProduction()
	e.ClientId, _ = global.Cfg.Int64("tencent_ad_client_id")
	e.ClientSecret = global.Cfg.String("tencent_ad_secret")
	e.GrantType = "authorization_code"
	e.OauthTokenOpts = &api.OauthTokenOpts{
		AuthorizationCode: optional.NewString(GetAuthorizationCode()),
	}
}

func (e *AccessToken) Run() string {
	tads := e.TAds
	ctx := *tads.Ctx
	response, _, err := tads.Oauth().Token(ctx, e.ClientId, e.ClientSecret, e.GrantType, e.OauthTokenOpts)
	if err != nil {
		if resErr, ok := err.(errors.ResponseError); ok {
			errStr, _ := json.Marshal(resErr)
			fmt.Println("Resopnse error:", string(errStr))
		} else {
			fmt.Println("Error:", err)
		}
	}
	return *response.AccessToken
}

func GetAuthorizationCode() (code string) {
	_, _ = util.HTTPGet(fmt.Sprintf("https://developers.e.qq.com/oauth/authorize?client_id=%s&redirect_uri=%s", global.Cfg.String("28072"), global.Cfg.String("tencent_ad_url")))
	cacheKey := cache.GetCacheKey(cache.TencentAdAuthorizationCode)
	var flag int
	for {
		flag++
		if flag > 5 {
			break
		}
		code = global.Cache.Get(cacheKey)
		time.Sleep(1)
	}
	return
}

func GetAccessToken() (token string) {
	cacheKey := cache.GetCacheKey(cache.TencentAdAccessToken)
	token = global.Cache.Get(cacheKey)
	if token == "" {
		e := &AccessToken{}
		e.Init()
		token = e.Run()
		if err := global.Cache.Set(cacheKey, token, 86400); err != nil {
			println("tencent_ad_token_err", err.Error())
		}
	}
	return
}
