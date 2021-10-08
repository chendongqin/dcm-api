package tencent_ad

import (
	"dongchamao/global"
	"dongchamao/global/utils"
	"encoding/json"
	"fmt"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
	"net/http"
)

type UserActionSets struct {
	TAds        *ads.SDKClient
	AccessToken string
	Data        model.UserActionSetsAddRequest
}

func (e *UserActionSets) Init(name string) {
	accountId := utils.ToInt64(global.Cfg.String("tencent_ad_account_id"))
	e.AccessToken = GetAccessToken()
	e.TAds = ads.Init(&config.SDKConfig{
		AccessToken: e.AccessToken,
		IsDebug:     true,
	})
	e.Data = model.UserActionSetsAddRequest{
		AccountId: &accountId,
		Name:      &name,
		Type_:     model.AmUserActionSetType_WEB,
	}
}

func (e *UserActionSets) Run() (model.UserActionSetsAddResponseData, http.Header, error) {
	tads := e.TAds
	ctx := *tads.Ctx
	return tads.UserActionSets().Add(ctx, e.Data)
}

func AddUserActionSets(Name string) {
	e := &UserActionSets{}
	e.Init(Name)
	response, headers, err := e.Run()
	if err != nil {
		if resErr, ok := err.(errors.ResponseError); ok {
			errStr, _ := json.Marshal(resErr)
			fmt.Println("Response error:", string(errStr))
		} else {
			fmt.Println("Error:", err)
		}
	}
	fmt.Println("Response data:", response)
	fmt.Println("Headers:", headers)
}
