package business

import (
	"dongchamao/global"
	"dongchamao/models/dcm"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	"github.com/silenceper/wechat/v2/officialaccount/user"
	"time"
)

type WechatBusiness struct {
}

func NewWechatBusiness() *WechatBusiness {
	return new(WechatBusiness)
}

func (receiver *WechatBusiness) CreateTempQrCode(scene string) (string, error) {
	tmpReq := basic.NewTmpQrRequest(1*3600*time.Second, scene)
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

func (receiver *WechatBusiness) SubscribeOfficial(userWechat *user.Info) error {
	dbSession := dcm.GetDbSession()
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
	if exist {
		_, err = dbSession.InsertOne(wechatModel)
	} else {
		_, err = dbSession.Where("unionid = ?", userWechat.UnionID).Cols("openid", "unionid", "nick_name", "avatar", "sex", "country", "province", "city", "language", "remark", "subscribe", "subscribe_time", "subscribe_scene").Update(wechatModel)
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
