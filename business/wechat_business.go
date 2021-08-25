package business

import (
	"dongchamao/global"
	"fmt"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	"time"
)

type WechatBusiness struct {
}

func NewWechatBusiness() *WechatBusiness {
	return new(WechatBusiness)
}

func (this *WechatBusiness) CreateTempQrCode(scene string) (string, error) {
	tmpReq := basic.NewTmpQrRequest(1*3600*time.Second, scene)
	tk, err := global.WxOfficial.GetBasic().GetQRTicket(tmpReq)
	fmt.Println(global.WxOfficial.GetAccessToken())
	if err != nil {
		return "", err
	}
	if tk.Ticket == "" {
		return "", errors.New("empty ticket")
	}
	return basic.ShowQRCode(tk), nil
}
