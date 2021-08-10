package controllers

import (
	"dongchamao/services/payer"
	"fmt"
)

type CallbackController struct {
	ApiBaseController
}

func (receiver *CallbackController) WechatNotify() {
	info, content, err := payer.Notify(receiver.Ctx.Request)
	fmt.Println("回调测试", info, content, err)
	receiver.SuccReturn(nil)
	return

}
