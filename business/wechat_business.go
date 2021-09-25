package business

import (
	"dongchamao/global"
	"dongchamao/models/dcm"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/menu"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/silenceper/wechat/v2/officialaccount/user"
	"time"
)

type WechatBusiness struct {
}

const (
	WechatMsgTemplateLiveMonitorFinish = "enHnUWWRMZiOdiyTlP4EKHy6sO6AX3SMlmPyamKNa3c" //直播监控结束通知
	WechatMsgTemplateLiveMonitorBegin  = "0JlZIQz3gk0hCYwNLtBmtd-WAGZ0zK_uQ_aR6-1ZpS0" //直播监控开始
	WechatMsgTemplateBindNotice        = "-O4v3TG7c8hV1xeVwt5jUdHLP0bLKghvfsxWMUCPeU8" //绑定通知 - 模板id
	WechatMsgTemplateLoginNotice       = "lbk1OMOLE1x-tSQIqap8otYWlAu-7d8mTmTnT0R4k50" //登录成功通知-模板id
	WechatMsgTemplateLoginOutNotice    = "e-tT-wntWlj-tVpb2YktMaPN5CnkBdKxGK11lq_ccUQ" //退出成功通知-模板id
	WechatMsgTemplateAmountExpire      = "e-tT-wntWlj-tVpb2YktMaPN5CnkBdKxGK11lq_ccUQ" //会员到期通知-模板id

)

func NewWechatBusiness() *WechatBusiness {
	return new(WechatBusiness)
}

func (receiver *WechatBusiness) CreateTempQrCode(scene string) (string, error) {
	tmpReq := basic.NewTmpQrRequest(3500, scene)
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

//获取用户微信基本信息
func (receiver *WechatBusiness) SendMsg(openId, templateID string, Data map[string]*message.TemplateDataItem, url string) error {
	msg := &message.TemplateMessage{
		ToUser:     openId,
		TemplateID: templateID, // 必须, 模版ID
		URL:        url,
		Data:       Data,
	}
	_, err := global.WxOfficial.GetTemplate().Send(msg)
	return err
}

func (receiver *WechatBusiness) SubscribeOfficial(userWechat *user.Info) error {
	dbSession := dcm.GetDbSession()
	dbSession.Begin()
	defer dbSession.Close()
	wechatModel := dcm.DcWechat{} //unionId 为主...
	exist, err := dbSession.Where("unionid = ?", userWechat.UnionID).Get(&wechatModel)
	wechatModel.Unionid = userWechat.UnionID
	wechatModel.Avatar = userWechat.Headimgurl
	wechatModel.NickName = userWechat.Nickname
	wechatModel.Openid = userWechat.OpenID
	wechatModel.Avatar = userWechat.Headimgurl
	wechatModel.Sex = int(userWechat.Sex)
	wechatModel.Country = userWechat.Country
	wechatModel.Province = userWechat.Province
	wechatModel.City = userWechat.City
	wechatModel.Groupid = int(userWechat.GroupID)
	wechatModel.Language = userWechat.Language
	wechatModel.Remark = userWechat.Remark
	wechatModel.Subscribe = int(userWechat.Subscribe)
	wechatModel.SubscribeTime = int64(userWechat.SubscribeTime)
	//wechatModel.UnsubscribeTime = 0
	wechatModel.SubscribeScene = userWechat.SubscribeScene
	wechatModel.QrScene = userWechat.QrScene
	wechatModel.QrSceneStr = userWechat.QrSceneStr
	if !exist {
		wechatModel.CreatedAt = time.Now()
		_, err = dbSession.InsertOne(wechatModel)
	} else {
		_, err = dbSession.Where("unionid = ?", userWechat.UnionID).Cols("openid", "unionid", "nick_name", "avatar",
			"sex", "country", "province", "city", "language", "remark", "subscribe", "subscribe_time", "subscribe_scene").
			Update(wechatModel)
	}
	if err != nil {
		dbSession.Rollback()
		return err
	}
	//填充user表的openid
	if userWechat.UnionID != "" && userWechat.OpenID != "" {
		_, _ = dbSession.Where("unionid = ?", userWechat.UnionID).Cols("openid").Update(&dcm.DcUser{Openid: userWechat.OpenID})
	}
	_ = dbSession.Commit()
	return nil

}

func (receiver *WechatBusiness) UnSubscribeOfficial(unionId string) error {
	_, err := dcm.GetDbSession().Table(dcm.DcWechat{}).Where("unionid = ?", unionId).
		Update(map[string]interface{}{"subscribe": 2, "unsubscribe_time": time.Now().Unix()})
	if err != nil {
		return err
	}
	return nil
}

func (receiver *WechatBusiness) SubscribeApp(userWechat *user.Info) error {
	dbSession := dcm.GetDbSession()
	dbSession.Begin()
	defer dbSession.Close()
	wechatModel := dcm.DcWechat{} //unionId 为主...
	exist, err := dbSession.Where("unionid = ?", userWechat.UnionID).Get(&wechatModel)
	wechatModel.Unionid = userWechat.UnionID
	wechatModel.NickName = userWechat.Nickname
	wechatModel.Avatar = userWechat.Headimgurl
	wechatModel.Sex = int(userWechat.Sex)
	wechatModel.Country = userWechat.Country
	wechatModel.Province = userWechat.Province
	wechatModel.City = userWechat.City
	wechatModel.Language = userWechat.Language
	wechatModel.Remark = userWechat.Remark
	wechatModel.Subscribe = int(userWechat.Subscribe)
	wechatModel.SubscribeTime = int64(userWechat.SubscribeTime)
	//wechatModel.UnsubscribeTime = 0
	wechatModel.SubscribeScene = userWechat.SubscribeScene
	wechatModel.QrScene = userWechat.QrScene
	wechatModel.QrSceneStr = userWechat.QrSceneStr
	wechatModel.Groupid = int(userWechat.GroupID)
	wechatModel.OpenidApp = userWechat.OpenID

	//如果不存在则添加，存在则更新
	if !exist {
		wechatModel.CreatedAt = time.Now()
		_, err = dbSession.InsertOne(wechatModel)
	} else {
		_, err = dbSession.Where("unionid = ?", userWechat.UnionID).Cols("openid_app", "unionid", "nick_name", "avatar",
			"sex", "country", "province", "city", "language", "remark", "subscribe", "subscribe_time", "subscribe_scene").
			Update(wechatModel)
	}
	if err != nil {
		dbSession.Rollback()
		return err
	}
	//填充user表的openid
	if userWechat.UnionID != "" && userWechat.OpenID != "" {
		_, _ = dbSession.Where("unionid = ?", userWechat.UnionID).Cols("openid_app").Update(&dcm.DcUser{OpenidApp: userWechat.OpenID})
	}
	_ = dbSession.Commit()
	return nil
}

func (receiver *WechatBusiness) BindWechat(userId int64, unionId string) {}

//菜单处理
func (receiver *WechatBusiness) GetMenus() (resMenu menu.ResMenu, err error) {
	return global.WxOfficial.GetMenu().GetMenu()
}

func (receiver *WechatBusiness) UpdateMenus(menuMap map[string]interface{}) error {
	menuByte, _ := jsoniter.Marshal(menuMap)
	menuOption := global.WxOfficial.GetMenu()
	return menuOption.SetMenuByJSON(string(menuByte))
}

//素材处理
func (receiver *WechatBusiness) GetMediaList(mediaType material.PermanentMaterialType, from, to int64) material.ArticleList {
	res, _ := global.WxOfficial.GetMaterial().BatchGetMaterial(mediaType, from, to)
	return res
}

func (receiver *WechatBusiness) AddMedia(mediaType material.MediaType, filename string) (string, string, error) {
	return global.WxOfficial.GetMaterial().AddMaterial(mediaType, filename)
}

func (receiver *WechatBusiness) DelMedia(mediaId string) error {
	return global.WxOfficial.GetMaterial().DeleteMaterial(mediaId)
}

func (receiver *WechatBusiness) GetMenuClick(key string) (click dcm.DcWechatMenuClick, err error) {
	_, err = dcm.GetDbSession().Table(dcm.DcWechatMenuClick{}).Where("`msg_key`=?", key).Get(&click)
	return
}

//绑定成功--发送微信模板消息
func (receiver *WechatBusiness) BindSendWechatMsg(user *dcm.DcUser) {
	if user.Openid == "" {
		return
	}
	msgMap := map[string]*message.TemplateDataItem{
		"first": {
			Value: "账号绑定成功",
			Color: "red",
		},
		"keyword1": {
			Value: user.Username,
			Color: "",
		},
		"keyword2": {
			Value: user.Nickname,
			Color: "",
		},
		"remark": {
			Value: "",
			Color: "",
		},
	}
	err := NewWechatBusiness().SendMsg(user.Openid, WechatMsgTemplateBindNotice, msgMap, DyDcmUrl)
	if err != nil {
		return
	}
	return
}

//登录成功--发送微信模板消息
func (receiver *WechatBusiness) LoginWechatMsg(user *dcm.DcUser) {
	if user.Openid == "" {
		return
	}
	msgMap := map[string]*message.TemplateDataItem{
		"first": {
			Value: "登录成功",
			Color: "red",
		},
		"keyword1": {
			Value: user.UpdateTime.Format("20060102"),
			Color: "",
		},
		"keyword2": {
			Value: user.LoginIp,
			Color: "",
		},
		"remark": {
			Value: "",
			Color: "",
		},
	}
	err := NewWechatBusiness().SendMsg(user.Openid, WechatMsgTemplateLoginNotice, msgMap, DyDcmUrl)
	if err != nil {
		return
	}
	return
}

//退出登录--发送微信模板消息
func (receiver *WechatBusiness) LoginOutWechatMsg(user *dcm.DcUser) {
	if user.Openid == "" {
		return
	}
	msgMap := map[string]*message.TemplateDataItem{
		"first": {
			Value: "退出登录",
			Color: "red",
		},
		"keyword1": {
			Value: user.Nickname,
			Color: "",
		},
		"keyword2": {
			Value: time.Now().Format("20060102"),
			Color: "",
		},
	}
	err := NewWechatBusiness().SendMsg(user.Openid, WechatMsgTemplateLoginOutNotice, msgMap, DyDcmUrl)
	if err != nil {
		return
	}
	return
}

func (receiver *WechatBusiness) AmountExpireWechatNotice(user *dcm.UserVipJpinCombine) {
	fmt.Printf("%+v\n", user)
	if user.DcJoinUser.Openid == "" {
		return
	}
	//DcUserVip
	vipBusiness := NewVipBusiness()
	levels := vipBusiness.GetUserLevels()

	msgMap := map[string]*message.TemplateDataItem{
		"first": {
			Value: "会员即将到期",
			Color: "red",
		},
		"keyword1": {
			Value: "会员等级",
			Color: "",
		},
		"keyword2": {
			Value: levels[user.DcJoinUserVip.Level],
			Color: "",
		},
		"keyword3": {
			Value: user.DcJoinUserVip.Expiration,
			Color: "",
		},
		"keyword4": {
			Value: user.DcJoinUser.Username,
			Color: "",
		},
	}
	err := NewWechatBusiness().SendMsg(user.Openid, WechatMsgTemplateAmountExpire, msgMap, DyDcmUrl)
	if err != nil {
		return
	}
	return
}
