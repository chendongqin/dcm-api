package tencent_ad

import (
	"dongchamao/global"
	"encoding/json"
	"fmt"
	"github.com/antihax/optional"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/api"
	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
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
	e.ClientId, _ = global.Cfg.Int64("1112002873")
	e.ClientSecret = global.Cfg.String("vFQvPWDjwOy4faq0")
	e.GrantType = "authorization_code"
	e.OauthTokenOpts = &api.OauthTokenOpts{
		AuthorizationCode: optional.NewString("YOUR AUTHORIZATION CODE"),
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

func GetAccessToken() string {
	e := &AccessToken{}
	e.Init()
	return e.Run()
}
