package controllers

import (
	"dongchamao/business"
	"dongchamao/global"
	"fmt"
)

type SpiderAuthController struct {
	ApiBaseController
}

func (this *SpiderAuthController) Prepare() {
	this.InitApi()
}

//====== 爬虫收集抖音号授权 ==============================

func (this *SpiderAuthController) GetQrCodeMcn() {
	s := business.NewDySpiderAuthScan()
	res, csrfToken, codeIP := s.GetQrCodeMcn(this.Ip)
	if res == nil {
		this.FailReturn(global.NewMsgError("获取二维码失败"))
		return
	}

	ret := new(business.SpiderAuthData)
	ret.CsrfToken = csrfToken
	ret.QrcodeIndexUrl, _ = res.Get("data").Get("qrcode_index_url").String()
	ret.Qrcode, _ = res.Get("data").Get("qrcode").String()
	ret.Token, _ = res.Get("data").Get("token").String()
	ret.CodeIP = codeIP

	this.SuccReturn(ret)
}

//扫完码  用token获取用户信息
func (this *SpiderAuthController) CheckQrConnectMcn() {
	InputDatas := this.InputFormat()
	token := InputDatas.GetString("token", "")
	csrfToken := InputDatas.GetString("csrf_token", "")
	codeIP := InputDatas.GetString("code_ip", "")

	success, cookies := business.NewDySpiderAuthScan().CheckQrConnectMcn(token, csrfToken, codeIP)
	if success == false {
		this.FailReturn(global.NewMsgError("绑定失败"))
		return
	}

	userInfo, _ := business.NewDySpiderAuthScan().SetCookie(cookies).GetUserInfo()
	fmt.Println(userInfo)
	//
	//auth := douyinmodelsV2.NewSvDySpiderOauth()
	//auth.DyUid = userInfo.Uid
	//auth.NickName = userInfo.Nickname
	//auth.AvatarUrl = userInfo.AvatarThumb
	//auth.ShortId = userInfo.ShortId
	//auth.NickName = userInfo.Nickname
	//auth.Cookies = douyinmodelsV2.NewSvDyCreatorOauthModel().CookieToString(cookies)
	//auth.SessionId = douyinmodelsV2.NewSvDySpiderOauth().GetSessionId(cookies)
	//exist := auth.AddOrUpdate()
	//if exist {
	//	this.FailReturn(global.NewMsgError("用户已绑定过"))
	//}

	this.SuccReturn(userInfo)
}
