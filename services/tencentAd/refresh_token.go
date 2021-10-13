package tencent_ad

import (
	"dongchamao/global"
	"dongchamao/global/cache"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"log"
)

type RefreshAccessToken struct {
	TAds           *ads.SDKClient
	AccessToken    string
	ClientId       int64
	ClientSecret   string
	GrantType      string
	OauthTokenOpts *api.OauthTokenOpts
}

func (e *RefreshAccessToken) Init(freshToken string) {
	e.TAds = ads.Init(&config.SDKConfig{})
	e.TAds.UseProduction()
	e.ClientId, _ = global.Cfg.Int64("tencent_ad_client_id")
	e.ClientSecret = global.Cfg.String("tencent_ad_secret")
	e.GrantType = "refresh_token"
	e.OauthTokenOpts = &api.OauthTokenOpts{
		RefreshToken: optional.NewString(freshToken),
	}
}

func (e *RefreshAccessToken) Run() {
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
	} else {
		tads.SetAccessToken(*response.AccessToken)
		cacheKey := cache.GetCacheKey(cache.TencentAdAccessToken)
		if err := global.Cache.Set(cacheKey, *response.AccessToken, 86400); err != nil {
			log.Println("tencent_ad_fresh_token_err:", err.Error())
		}
	}
}

func GetAccessToken() (token string) {
	cacheKey := cache.GetCacheKey(cache.TencentAdAccessToken)
	token = global.Cache.Get(cacheKey)
	if token == "" {
		freshCacheKey := cache.GetCacheKey(cache.TencentAdRefreshToken)
		freshToken := global.Cache.Get(freshCacheKey)
		r := &RefreshAccessToken{}
		r.Init(freshToken)
		r.Run()
		token = r.AccessToken
		if err := global.Cache.Set(cacheKey, token, 86400); err != nil {
			log.Println("tencent_ad_token_err", err.Error())
		}
	}
	return
}
