package business

import (
	"dongchamao/global"
	"dongchamao/models/dcm"
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
func (receiver *WechatBusiness) SendMsg(openId, templateID, url string, Data map[string]*message.TemplateDataItem) error {
	//map[string]*message.TemplateDataItem{
	//	"first": {
	//		Value: "如故，欢迎您~",
	//		Color: "",
	//	},
	//	"keyword1": {
	//		Value: time.Now().Format("2006-01-02 15:04:05"),
	//		Color: "",
	//	},
	//	"keyword2": {
	//		Value: "127.0.0.1",
	//		Color: "",
	//	},
	//	"remark": {
	//		Value: "如有非本人操作，请及时修改密码~",
	//		Color: "red",
	//	},
	//},
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
	return global.WxOfficial.GetMenu().SetMenuByJSON(string(menuByte))
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
