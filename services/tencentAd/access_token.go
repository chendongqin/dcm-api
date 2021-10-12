package tencent_ad

import (
	"dongchamao/business"
	"dongchamao/global"
	"dongchamao/global/cache"
	"encoding/json"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"log"
)

type AccessToken struct {
	TAds           *ads.SDKClient
	AccessToken    string
	ClientId       int64
	ClientSecret   string
	GrantType      string
	OauthTokenOpts *api.OauthTokenOpts
}

func (e *AccessToken) Init(authorizationCode string) {
	e.TAds = ads.Init(&config.SDKConfig{})
	e.TAds.UseProduction()
	e.ClientId, _ = global.Cfg.Int64("tencent_ad_client_id")
	e.ClientSecret = global.Cfg.String("tencent_ad_secret")
	e.GrantType = "authorization_code"
	e.OauthTokenOpts = &api.OauthTokenOpts{
		AuthorizationCode: optional.NewString(authorizationCode),
		RedirectUri:       optional.NewString(global.Cfg.String("tencent_ad_url")),
	}
}

func (e *AccessToken) Run() {
	tads := e.TAds
	ctx := *tads.Ctx
	response, _, err := tads.Oauth().Token(ctx, e.ClientId, e.ClientSecret, e.GrantType, e.OauthTokenOpts)
	if err != nil {
		if resErr, ok := err.(errors.ResponseError); ok {
			errStr, _ := json.Marshal(resErr)
			log.Println("Resopnse error:", string(errStr))
		} else {
			log.Println("Error:", err)
		}
	}
	freshCacheKey := cache.GetCacheKey(cache.TencentAdRefreshToken)
	cacheKey := cache.GetCacheKey(cache.TencentAdAccessToken)
	business.NewMonitorBusiness().SendErr("AccessToken:", *response.AccessToken+"==="+*response.RefreshToken)
	if err := global.Cache.Set(cacheKey, *response.AccessToken, 86400); err != nil {
		business.NewMonitorBusiness().SendErr("tencent_ad_fresh_token_err:", err.Error())
	}
	if err := global.Cache.Set(freshCacheKey, *response.RefreshToken, 999999999); err != nil {
		business.NewMonitorBusiness().SendErr("tencent_ad_token_err:", err.Error())
	}
	return
}
