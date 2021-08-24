package v1

import (
	"dongchamao/controllers/api"
	"github.com/prometheus/common/log"
)

type WechatController struct {
	controllers.ApiBaseController
}

func (receiver *WechatController) Receive() {
	InputData := receiver.InputFormat()
	log.Info("微信回调数据:", InputData)
}
