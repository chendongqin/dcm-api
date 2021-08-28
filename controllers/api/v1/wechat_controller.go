package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/consistent"
	"dongchamao/global/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"strings"
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
	cacheKey := "openid:" + sessionId

	openId := global.Cache.Get(cacheKey)
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
		text = message.NewText("扫码登录成功!!") //TODO 事件推送返回信息 可以抽象出来 也可以后台配置
		//msg.EventKey 返回场景基本都是qrscene_你自己定义场景key
		if msg.MsgType == message.MsgTypeEvent {
			switch msg.Event {
			case message.EventSubscribe:
				//自定义evnet_key 解析
				logs.Error("[扫码登录微信1001]=>缓存key:[%s],openid:[%s]", msg.EventKey, msg.GetOpenID())
				if strings.Contains(msg.EventKey, consistent.WECHAT_QR_LOGIN) {
					//这边自定义的事件 还会有细分  qrlogin / qrscene_qrlogin
					if strings.Contains(msg.EventKey, consistent.WECHAT_QR_SCENE_LOGIN) {
						msg.EventKey = msg.EventKey[8:]
					}
					err := business.NewWechatBusiness().SubscribeOfficial(userWechat)
					if err != nil {
						logs.Error("[扫码绑定] 数据更新失败1001，err: %s", err)
						text = message.NewText("扫码关注失败!!")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					//设置 openid 缓存 前端监听
					_ = global.Cache.Set("openid:"+msg.EventKey, msg.GetOpenID(), 1800)
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			case message.EventUnsubscribe:
				logs.Error("[扫码登录微信1002]=>缓存key:[%s],openid:[%s]", msg.EventKey, msg.GetOpenID())
				_ = business.NewWechatBusiness().UnSubscribeOfficial(msg.GetOpenID())
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			case message.EventScan:
				//自定义事件key
				if strings.Contains(msg.EventKey, consistent.WECHAT_QR_LOGIN) {
					err := business.NewWechatBusiness().SubscribeOfficial(userWechat)
					if err != nil {
						logs.Error("[扫码绑定] 数据更新失败1001，err: %s", err)
						text = message.NewText("扫码关注失败!!")
						return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
					}
					//设置 openid 缓存 前端监听
					_ = global.Cache.Set("openid:"+msg.EventKey, msg.GetOpenID(), 1800)
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				//default:
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
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

//微信客户端 获取open相关信息

//func (receiver *WechatController) WxApp() {
//	inputData := receiver.InputFormat()
//	code := inputData.GetString("code", "")
//}
