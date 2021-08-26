package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
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

//扫码微信回调成功 通知前端
func (receiver *WechatController) CheckScan() {
	inputData := receiver.InputFormat()
	sessionId := inputData.GetString("session_id", "")
	if sessionId == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	//从缓存中获取用户openId
	openId := global.Cache.Get("openid:" + sessionId)
	if openId == "" {
		receiver.FailReturn(global.NewError(4006))
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"open_id": openId,
	})
	return
}

func (receiver *WechatController) Receive() {
	inputData := receiver.InputFormat()
	server := global.WxOfficial.GetServer(receiver.Ctx.Request, receiver.Ctx.ResponseWriter)
	if beego.BConfig.RunMode == beego.DEV { //测试环境默认通过
		server.SkipValidate(true)
	}
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		logs.Error("[微信回调]=>请求参数:[%s],事件内容:[%s]", inputData, msg)
		//回复消息：演示回复用户发送的消息
		userWechat, err := business.NewWechatBusiness().GetInfoByOpenId(msg.GetOpenID())
		if err != nil {
			logs.Error("[微信回调] 获取用户信息失败, err: %s", err)
			return nil
		}
		var text *message.Text
		text = message.NewText("扫码登录成功!!")
		if msg.MsgType == message.MsgTypeEvent {
			switch msg.Event {
			case message.EventSubscribe:
				logs.Error("[扫码登录微信1001]=>缓存key:[%s],openid:[%s]", msg.EventKey, msg.GetOpenID())
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			case message.EventUnsubscribe:
				logs.Error("[扫码登录微信1002]=>缓存key:[%s],openid:[%s]", msg.EventKey, msg.GetOpenID())
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			case message.EventScan:
				//自定义事件key
				err := business.NewWechatBusiness().SubscribeOfficial(userWechat)
				if err != nil {
					logs.Error("[扫码绑定] 数据更新失败1001，err: %s", err)
					return nil
				}
				text = message.NewText("扫码登录成功！")
				//设置 openid 缓存 前端监听
				logs.Error("[扫码登录微信1003]=>缓存key:[%s],openid:[%s]", msg.EventKey, msg.GetOpenID())
				_ = global.Cache.Set("openid:"+msg.EventKey, msg.GetOpenID(), 1800)
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				//default:
				//	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			}
		}
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		logs.Error("[微信回复] 回复消息失败 serve err: %s", err.Error())
		return
	}
	//发送回复的消息
	_ = server.Send()
}
