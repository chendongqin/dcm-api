package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/common/log"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type WechatController struct {
	controllers.ApiBaseController
}

//临时二维码...
func (receiver *WechatController) QrCode() {
	inputData := receiver.InputFormat()
	sessionId := inputData.GetString("session_id", "")
	var err error
	var qrUrl string
	if sessionId != "" {
		qrUrl = global.Cache.Get("qrcode:" + sessionId)
	}
	if qrUrl == "" {
		sessionId = "qrlogin:" + utils.Md5_encode(utils.Date("", 0)+"|"+utils.GetRandomStringNew(12))
		qrUrl, err = business.NewWechatBusiness().CreateTempQrCode(sessionId)
		if err == nil && qrUrl != "" {
			_ = global.Cache.Set("qrcode:"+sessionId, qrUrl, 300)
		} else {
			receiver.FailReturn(global.NewCommonError(err))
			return
		}
	}
	receiver.SuccReturn(map[string]interface{}{
		"qr_code":    qrUrl,
		"session_id": sessionId,
	})
}

func (receiver *WechatController) Receive() {
	//wxBusiness := business.NewWechatBusiness()
	InputData := receiver.InputFormat()
	log.Info("微信回调数据:", InputData)
	server := global.WxOfficial.GetServer(receiver.Ctx.Request, receiver.Ctx.ResponseWriter)
	if beego.BConfig.RunMode == beego.DEV { //测试环境默认通过
		server.SkipValidate(true)
	}
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		logs.Error("微信调数据msg：", msg)
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		logs.Error("[微信回复] 回复消息失败 serve err: %s", err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		logs.Error("[微信回复] 回复消息失败 send err: %s", err)
		return
	}
}
