package controllers

import (
	"dongchamao/business"
	"dongchamao/global"
	"dongchamao/models/dcm"
	"time"
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
	dbSession := dcm.GetDbSession()
	dySpiderAuthScan := business.NewDySpiderAuthScan()
	defer dbSession.Close()
	token := InputDatas.GetString("token", "")
	csrfToken := InputDatas.GetString("csrf_token", "")
	codeIP := InputDatas.GetString("code_ip", "")

	success, cookies := business.NewDySpiderAuthScan().CheckQrConnectMcn(token, csrfToken, codeIP)
	if success == false {
		this.FailReturn(global.NewError(6102))
		return
	}

	userInfo, _ := dySpiderAuthScan.SetCookie(cookies).GetUserInfo()
	auth := &dcm.DySpiderAuth{}
	auth.Uid = userInfo.Uid
	auth.Nickname = userInfo.Nickname
	auth.Cookies = dySpiderAuthScan.CookieToString(cookies)
	auth.Sessionid = dySpiderAuthScan.GetSessionId(cookies)
	exist, _ := dbSession.Table(auth).Where("uid = ?", auth.Uid).Exist()
	if exist == false {
		auth.CreateTime = time.Now()
		_, _ = dbSession.Insert(auth)
	} else {
		auth.UpdateTime = time.Now()
		_, _ = dbSession.Where("uid = ?", auth.Uid).Update(auth)
		this.FailReturn(global.NewError(6101))
	}
	this.SuccReturn("ok")

}
