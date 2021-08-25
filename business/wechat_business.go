package business

import (
	"dongchamao/global"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	"github.com/silenceper/wechat/v2/officialaccount/user"
	"time"
)

type WechatBusiness struct {
}

func NewWechatBusiness() *WechatBusiness {
	return new(WechatBusiness)
}

func (receiver *WechatBusiness) CreateTempQrCode(scene string) (string, error) {
	tmpReq := basic.NewTmpQrRequest(1*3600*time.Second, scene)
	tk, err := global.WxOfficial.GetBasic().GetQRTicket(tmpReq)
	if err != nil {
		return "", err
	}
	if tk.Ticket == "" {
		return "", errors.New("empty ticket")
	}
	return basic.ShowQRCode(tk), nil
}

//获取用户微信基本信息
func (receiver *WechatBusiness) GetInfoByOpenId(openId string) (*user.Info, error) {
	userWechat, err := global.WxOfficial.GetUser().GetUserInfo(openId)
	if err != nil {
		return nil, err
	} else {
		return userWechat, nil
	}
}
