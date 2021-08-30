package business

import (
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/officialaccount/user"
	"net/url"
	"time"
)

const (
	AccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	// RefreshTokenURL 重新获取access_token
	RefreshTokenURL = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	// UserInfoURL 通过access_token获取userInfo
	UserInfoURL = "https://api.weixin.qq.com/sns/userinfo"
)

type WxAppBusiness struct {
	AppID  string `json:"appid"`  // 微信APPID
	Secret string `json:"secret"` // 微信Secret
}

func NewWxAppBusiness() *WxAppBusiness {
	wxApp := new(WxAppBusiness)
	wxApp.AppID = "wxe57558084988e26d"
	wxApp.Secret = "d1151b4ebfc08503ae0d2b9ea4f290d9"
	return wxApp
}

// WxAccessToken 微信授权Token
type WxAccessToken struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    uint   `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	OpenID       string `json:"openid,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ErrCode      uint   `json:"errcode,omitempty"`
	ErrMsg       string `json:"errmsg,omitempty"`
	ExpiredAt    time.Time
}

// AppLogin 微信APP登录 直接登录获取用户信息
func (receiver *WxAppBusiness) AppLogin(code string) (unionid string, err error) {
	userInfo, err := receiver.LoginCode(code)
	if err != nil {
		return "", err
	}
	err = NewWechatBusiness().SubscribeApp(userInfo)
	if err != nil {
		return "", err
	}

	return userInfo.UnionID, nil
}

// LoginCode 通过Code登录
func (receiver *WxAppBusiness) LoginCode(code string) (wxUserInfo *user.Info, err error) {
	accessToken, err := receiver.GetWxAccessToken(code)

	if err != nil {
		return wxUserInfo, err
	}
	return accessToken.GetUserInfo()
}

// GetWxAccessToken 通过code获取AccessToken
func (receiver *WxAppBusiness) GetWxAccessToken(code string) (accessToken *WxAccessToken, err error) {

	if code == "" {
		return accessToken, errors.New("GetWxAccessToken error: code is null")
	}

	params := url.Values{
		"code":       []string{code},
		"grant_type": []string{"authorization_code"},
	}

	t, err := utils.Struct2Map(receiver)
	if err != nil {
		return accessToken, err
	}

	for k, v := range t {
		params.Set(k, v)
	}

	body, err := utils.NewRequest("GET", AccessTokenURL, []byte(params.Encode()))
	if err != nil {
		return accessToken, err
	}

	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		return accessToken, err
	}
	if accessToken.ErrMsg != "" {
		return accessToken, errors.New(accessToken.ErrMsg)
	}

	return
}

// GetUserInfo 获取用户资料
func (receiver *WxAccessToken) GetUserInfo() (wxUserInfo *user.Info, err error) {
	if receiver.AccessToken == "" {
		return nil, errors.New("GetUserInfo error: accessToken is null")
	}

	if receiver.OpenID == "" {
		return nil, errors.New("GetUserInfo error: openID is null")
	}

	params := url.Values{
		"access_token": []string{receiver.AccessToken},
		"openid":       []string{receiver.OpenID},
	}
	body, err := utils.NewRequest("GET", UserInfoURL, []byte(params.Encode()))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &wxUserInfo)
	if err != nil {
		return nil, err
	}
	if wxUserInfo.OpenID == "" {
		return wxUserInfo, errors.New(wxUserInfo.ErrMsg)
	}

	return
}

func (receiver *WxAppBusiness) SubscribeApp(userWechat *user.Info) error {
	dbSession := dcm.GetDbSession()
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
		_, err = dbSession.Where("unionid = ?", userWechat.UnionID).Cols("openid", "unionid", "nick_name", "avatar",
			"sex", "country", "province", "city", "language", "remark", "subscribe", "subscribe_time", "subscribe_scene").
			Update(wechatModel)
	}
	if err != nil {
		return err
	}
	return nil
}

// GetRefreshToken 重新获取AccessToken

// CheckAccessToken 校验AccessToken
