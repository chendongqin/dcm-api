package controllers

import (
	"dongchamao/services/payer"
	"github.com/astaxie/beego/logs"
)

type CallbackController struct {
	ApiBaseController
}

func (receiver *CallbackController) WechatNotify() {
	info, content, err := payer.Notify(receiver.Ctx.Request)
	logs.Error("回调测试", info.Summary, info.RawRequest, err)
	logs.Error("回调饭回数据", content)
	receiver.SuccReturn(nil)
	return

}
