package v1

import (
	"dongchamao/business"
	"dongchamao/controllers/api"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/consistent"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"dongchamao/models/repost/dy"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"strings"
	"time"
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
	//从缓存中获取用户unionid
	cacheKey := "unionid:" + sessionId
	uniondId := global.Cache.Get(cacheKey)
	if uniondId == "" {
		receiver.FailReturn(global.NewError(4006))
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"unionid": business.IdEncrypt(uniondId),
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
		logs.Info("[微信回调]=>请求参数:[%s]", inputData)
		//回复消息：演示回复用户发送的消息
		openId := msg.GetOpenID()
		userWechat, err := business.NewWechatBusiness().GetInfoByOpenId(openId)
		logs.Info("openId:%+v,userWechat:%+v", openId, userWechat)
		if err != nil {
			logs.Error("[微信回调] 获取用户信息失败, err: %s", err)
			return nil
		}
		var text *message.Text
		text = message.NewText(global.WECHATLOGINMSG) //TODO 事件推送返回信息 可以抽象出来 也可以后台配置
		//msg.EventKey 返回场景基本都是qrscene_你自己定义场景key
		logs.Info("MsgType:%s,Event:%s,EventKey:%s", msg.MsgType, msg.Event, msg.EventKey)
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
					//设置openid缓存 前端监听
					_ = global.Cache.Set("unionid:"+msg.EventKey, userWechat.UnionID, 1800)
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
				}
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			case message.EventUnsubscribe:
				logs.Error("[扫码登录微信1002]=>缓存key:[%s],openid:[%s]", msg.EventKey, userWechat.UnionID)
				_ = business.NewWechatBusiness().UnSubscribeOfficial(userWechat.UnionID)
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
					_ = global.Cache.Set("unionid:"+msg.EventKey, userWechat.UnionID, 1800)
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("扫码登陆成功！！\n\n" + global.WECHATLOGINMSG)}
				}
			case message.EventClick:
				click, _ := business.NewWechatBusiness().GetMenuClick(msg.EventKey)
				logs.Error("msg.EventKey:%s,click.Key:%s  click.Type:%s", msg.EventKey, click.MsgKey, message.MsgType(click.Type))
				var msgData interface{}
				if click.MsgKey != "" {
					switch message.MsgType(click.Type) {
					case message.MsgTypeText:
						msgData = message.NewText(click.Content)
						break
					case message.MsgTypeImage:
						msgData = message.NewImage(click.MediaId)
						break
					case message.MsgTypeVoice:
						msgData = message.NewVoice(click.MediaId)
						break
					case message.MsgTypeVideo:
						msgData = message.NewVideo(click.MediaId, click.Title, click.Description)
						break
					}
				} else {
					click.Type = "text"
					msgData = message.NewText("消息未知")
				}
				return &message.Reply{MsgType: message.MsgType(click.Type), MsgData: msgData}
			default:
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: "消息未知"}
			}
		} else if msg.MsgType == message.MsgTypeText {
			click, _ := business.NewWechatBusiness().GetMenuClick(msg.Content)
			logs.Info("msg.EventKey:%s,click.Key:%s  click.Type:%s", msg.EventKey, click.MsgKey, message.MsgType(click.Type))
			var msgData interface{}
			if click.MsgKey != "" {
				switch message.MsgType(click.Type) {
				case message.MsgTypeText:
					msgData = message.NewText(click.Content)
					break
				case message.MsgTypeImage:
					msgData = message.NewImage(click.MediaId)
					break
				case message.MsgTypeVoice:
					msgData = message.NewVoice(click.MediaId)
					break
				case message.MsgTypeVideo:
					msgData = message.NewVideo(click.MediaId, click.Title, click.Description)
					break
				}
			} else {
				click.Type = "text"
				msgData = message.NewText("消息未知")
			}
			return &message.Reply{MsgType: message.MsgType(click.Type), MsgData: msgData}
		} else {
			var msgData = message.NewText("消息未知")
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: msgData}
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
func (receiver *WechatController) WechatApp() {
	inputData := receiver.InputFormat()
	code := inputData.GetString("code", "")
	if code == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	unionid, err := business.NewWxAppBusiness().AppLogin(code)
	if err != nil {
		receiver.FailReturn(global.NewError(4302))
		return
	}
	receiver.SuccReturn(map[string]interface{}{
		"unionid": business.IdEncrypt(unionid),
	})
	return
}

//微信登录绑定老账号
func (receiver *WechatController) WechatPhone() {
	inputData := receiver.InputFormat()
	userName := inputData.GetString("username", "")
	code := inputData.GetString("code", "")
	unionid := business.IdDecrypt(inputData.GetString("unionid", ""))
	//source := inputData.GetString("source", "") //    1.二维码2.微信一键登录
	if userName == "" || unionid == "" {
		receiver.FailReturn(global.NewError(4000))
		return
	}
	dbSession := dcm.GetDbSession()
	wechatModel := dcm.DcWechat{}
	if exist, _ := dbSession.Where("unionid = ?", unionid).Get(&wechatModel); !exist {
		receiver.FailReturn(global.NewError(4304))
		return
	}
	//手机验证
	codeKey := cache.GetCacheKey(cache.SmsCodeVerify, "bind_mobile", userName)
	verifyCode := global.Cache.Get(codeKey)
	if verifyCode != code {
		receiver.FailReturn(global.NewError(4209))
		return
	}
	//查询手机是否该用户 ...
	userModel := dcm.DcUser{}
	if exist, _ := dbSession.Where("username = ?", userName).Get(&userModel); !exist {
		userModel.Username = userName
		userModel.Nickname = wechatModel.NickName
		userModel.Salt = utils.GetRandomString(4)
		userModel.Password = utils.Md5_encode(utils.GetRandomString(16) + userModel.Salt)
		userModel.Status = 1
		userModel.CreateTime = time.Now()
		userModel.UpdateTime = time.Now()
		//来源
		userModel.Entrance = business.AppIdMap[receiver.AppId]
		userModel.Channel = receiver.Channel
		userModel.ChannelWords = receiver.ChannelWords
		affect, err := dcm.Insert(nil, &userModel)
		if affect == 0 || err != nil {
			receiver.FailReturn(global.NewError(5000))
			return
		}
		business.NewUserBusiness().SendUserVip(&userModel, 7)
	}
	//开始更新用户信息
	if userModel.Unionid != "" {
		receiver.FailReturn(global.NewError(4305))
		return
	}

	//userModel.OpenidApp = wechatModel.OpenidApp
	//userModel.Openid = wechatModel.Openid
	userBusiness := business.NewUserBusiness()
	updateData := map[string]interface{}{
		"openid_app":  wechatModel.OpenidApp,
		"openid":      wechatModel.Openid,
		"unionid":     wechatModel.Unionid,
		"nickname":    wechatModel.NickName,
		"avatar":      wechatModel.Avatar,
		"login_time":  utils.GetNowTimeStamp(),
		"login_ip":    receiver.Ip,
		"update_time": utils.GetNowTimeStamp(),
	}
	affect, _ := userBusiness.UpdateUserAndClearCache(nil, userModel.Id, updateData)
	if affect == 0 {
		receiver.FailReturn(global.NewError(4213))
		return
	}
	tokenString, expire, err := userBusiness.CreateToken(receiver.AppId, userModel.Id, 604800)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	err = userBusiness.AddOrUpdateUniqueToken(userModel.Id, receiver.AppId, tokenString)
	if err != nil {
		receiver.FailReturn(global.NewError(5000))
		return
	}
	receiver.RegisterLogin(tokenString, expire)

	//绑定手机成功通知
	business.NewWechatBusiness().LoginWechatMsg(&userModel)
	receiver.SuccReturn(map[string]interface{}{
		"token_info": dy.RepostAccountToken{
			UserId:      userModel.Id,
			TokenString: tokenString,
			ExpTime:     expire,
		},
	})
	return
}
