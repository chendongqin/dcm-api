package business

import (
	"dongchamao/global"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/qr"
	"time"
)

type WechatBusiness struct {
}

func NewWechatBusiness() *WechatBusiness {
	return new(WechatBusiness)
}

func (this *WechatBusiness) CreateTempQrCode(scene string) (string, error) {
	req := qr.NewTmpQrRequest(1*3600*time.Second, scene)
	tk, err := global.WechatInstance.GetQR().GetQRTicket(req)
	if err != nil {
		return "", err
	}
	if tk.Ticket == "" {
		return "", errors.New("empty ticket")
	}
	return qr.ShowQRCode(tk), nil
}
