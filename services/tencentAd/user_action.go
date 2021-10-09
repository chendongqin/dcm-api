package tencent_ad

import (
	"dongchamao/global"
	"encoding/json"
	"fmt"
	"github.com/tencentad/marketing-api-go-sdk/pkg/ads"
	"github.com/tencentad/marketing-api-go-sdk/pkg/config"
	"github.com/tencentad/marketing-api-go-sdk/pkg/errors"
	"github.com/tencentad/marketing-api-go-sdk/pkg/model"
	"net/http"
	"time"
)

type UserActions struct {
	TAds        *ads.SDKClient
	AccessToken string
	Data        model.UserActionsAddRequest
}

func (e *UserActions) Init(Channel string) {
	var (
		now = time.Now().Unix()
	)
	CustomAction := "REGISTER"
	UserActionSetId := int64(0)
	accountId, _ := global.Cfg.Int64("tencent_ad_account_id")
	e.AccessToken = GetAccessToken()
	e.TAds = ads.Init(&config.SDKConfig{
		AccessToken: e.AccessToken,
	})
	var channelVal model.ActionChannelType
	channelVal = model.ActionChannelType_NATURAL
	if Channel == "0024" {
		channelVal = model.ActionChannelType_TENCENT
	}
	e.Data = model.UserActionsAddRequest{
		AccountId: &accountId,
		Actions: &[]model.UserAction{{
			ActionTime:   &now,
			UserId:       &model.ActionsUserId{},
			Channel:      channelVal,
			ActionType:   model.ActionType_CUSTOM,
			CustomAction: &CustomAction,
		}},
		UserActionSetId: &UserActionSetId,
	}
}

func (e *UserActions) RunExample() (interface{}, http.Header, error) {
	tads := e.TAds
	ctx := *tads.Ctx
	return tads.UserActions().Add(ctx, e.Data)
}

func UserActionsAdd(Channel string) {
	e := &UserActions{}
	e.Init(Channel)
	response, headers, err := e.RunExample()
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
