package controllers

import (
	"bytes"
	"dongchamao/business"
	"dongchamao/controllers"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/dcm"
	aliLog "dongchamao/services/ali_log"
	"dongchamao/services/elasticsearch"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type ApiBaseController struct {
	GetDatas         map[string]interface{}
	ApiDatas         interface{}
	UserId           int
	UserInfo         dcm.DcUser
	IsMonitor        bool
	AppId            int
	DyLevel          int
	XhsLevel         int
	TbLevel          int
	Ip               string
	TrueUri          string
	IsInitToken      bool               //是否初始化过token
	LastInitTokenErr global.CommonError //记录首次初始化token的错误
	Token            string
	HasAuth          bool
	HasLogin         bool
	MaxTotal         int
	controllers.BaseController
}

func (this *ApiBaseController) Prepare() {
	this.InitApiController()
}

func (this *ApiBaseController) InitApiController() {
	this.InitApi()
	this.AsfCheck()
	this.CheckSign()
	this.InitUserToken()
	//todo 上线白名单过滤
	//if this.AppId < 20000 {
	if this.AppId <= 10000 {
		if !utils.InArrayString(this.TrueUri, []string{"/v1/user/login", "/v1/account/info", "/v1/config/list", "/v1/sms/verify", "/v1/sms/code",
			"/v1/wechat/phone",
			"/v1/account/logout", "/v1/wechat/check", "/v1/wechat/qrcode", "/v1/account/password", "/v1/callback/wechat"}) {
			if business.WitheUsername(this.UserInfo.Username) != nil {
				this.FailReturn(global.NewError(88888))
				return
			}
		}
	}
}

func (this *ApiBaseController) IsMobileRequest() (is bool, version string) {
	userAgent := this.Ctx.Request.Header.Get("User-Agent")
	if strings.Contains(userAgent, "dongchamao") {
		data := strings.Split(userAgent, "/")
		version := ""
		if len(data) > 1 {
			version = strings.Trim(data[1], " ")
		}
		return true, version
	}
	return false, ""
}

func (this *ApiBaseController) IsAndroid() (bool, string) {
	userAgent := this.Ctx.Request.Header.Get("User-Agent")
	if strings.Contains(userAgent, "dongchamao-android") {
		data := strings.Split(userAgent, "/")
		version := ""
		if len(data) > 1 {
			version = strings.Trim(data[1], " ")
		}
		return true, version
	}
	return false, ""
}

func (this *ApiBaseController) IsIOS() (bool, string) {
	userAgent := this.Ctx.Request.Header.Get("User-Agent")
	if strings.Contains(userAgent, "dongchamao-ios") {
		data := strings.Split(userAgent, "/")
		version := ""
		if len(data) > 1 {
			version = strings.Trim(data[1], " ")
		}
		return true, version
	}
	return false, ""
}

func (this *ApiBaseController) CheckIp() {
	//获得IP白名单，IP黑名单配置
	this.Ip = this.Ctx.Input.IP()
	//whitelists := apiv1models.NewConfigModel().GetConfig("ip_whitelists", true)
	//blacklists := apiv1models.NewConfigModel().GetConfig("ip_blacklists", true)

	//black := strings.Split(blacklists, ",")
	//white := strings.Split(whitelists, ",")
	//if len(black) > 0 {
	//	if utils.InArray(this.Ip, black) == true && utils.InArray(this.Ip, white) == false {
	//		//记录本次异常请求IP
	//		logs.Info("blackip:ban:(" + this.Ip + ")")
	//		this.FailReturn(global.NewError(40011))
	//		return
	//	}
	//}
}

func (this *ApiBaseController) InitUserToken() (commonErr global.CommonError) {
	//如果已经完成初始化，将上一次初始化的错误直接返回
	if this.IsInitToken {
		return this.LastInitTokenErr
	}
	defer func() {
		// 记录初始化状态
		this.IsInitToken = true
		// 记录初始化错误
		this.LastInitTokenErr = commonErr
	}()

	tokenString := this.Ctx.Input.Cookie(global.LOGINCOOKIENAME)
	//cookie没有身份信息  从头部获取
	if tokenString == "" {
		tokenString = strings.Replace(this.Ctx.Input.Header("Authorization"), "Bearer ", "", 1)
	}
	if tokenString == "" {
		return global.NewError(4001)
	}
	token, _ := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Cfg.String("auth_code")), nil
	})
	if token == nil {
		return global.NewError(4001)
	}
	userBusiness := business.NewUserBusiness()
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jsonstr, err0 := json.Marshal(claims)
		tokenStruct := &business.TokenData{}
		err0 = json.Unmarshal(jsonstr, tokenStruct)
		if err0 != nil {
			return global.NewError(4001)
		}
		expireTime := tokenStruct.ExpireTime
		if utils.Time() > expireTime {
			return global.NewError(4002)
		} else {
			this.UserId = tokenStruct.Id
		}
		if this.UserId == 0 {
			return global.NewError(4001)
		}
		this.AppId = tokenStruct.AppId
		this.Token = tokenString
		userInfo, exist := userBusiness.GetCacheUser(this.UserId, true)
		if !exist {
			return global.NewError(4001)
		}
		this.UserInfo = userInfo
		//判断用户状态x
		if this.UserInfo.Status == 0 {
			return global.NewError(4212)
		}
		//除bindphone外的接口 没有phone不让访问
		if this.UserInfo.Username == "" {
			return global.NewError(4005)
		}

		//处理连续登录次数统计处理
		if userBusiness.UpdateVisitedTimes(this.UserInfo) {
			//发生更新，才需要清除缓存
			userBusiness.DeleteUserInfoCache(this.UserInfo.Id)
		}

		//验证 user Platform token的唯一性
		uniqueToken, exist := userBusiness.GetUniqueToken(this.UserId, this.AppId, true)
		if exist == false {
			return global.NewError(4003)
		} else {
			if tokenString != uniqueToken {
				//cmmlog.LoginLog(this.Ctx.Input.Header("X-Request-Id"), this.Ctx.Input.Header("X-Client-Id"), this.appId, this.UserId, "current:"+tokenString+"|unique:"+uniqueToken, this.Ctx.Input.IP(), this.Ctx.Request.UserAgent(), "unique", "unique_loginout")
				this.RegisterLogout()
				return global.NewError(4001)
			}
		}

		if this.AppId == 10000 || this.AppId == 10001 {
			this.RegisterLogin(tokenString, expireTime)
		}

		this.DyLevel = userBusiness.GetCacheUserLevel(this.UserId, 1, true)
		this.XhsLevel = userBusiness.GetCacheUserLevel(this.UserId, 2, true)
		this.TbLevel = userBusiness.GetCacheUserLevel(this.UserId, 3, true)

		//异步记录用户行为日志
		this.LogInputOutput("Format", this.ApiDatas)
		this.HasLogin = true
	} else {
		return global.NewError(4001)
	}
	return
}

func (this *ApiBaseController) CheckSign() {
	authBusiness := business.NewAccountAuthBusiness()
	appId := this.Ctx.Input.Header("APPID")
	if appId == "" {
		appId = "10000"
	}
	this.AppId = utils.ToInt(appId)
	if utils.InArrayInt(this.AppId, []int{10000, 10001, 10002, 10003, 10004, 10005}) {
		if authBusiness.AuthSignWhiteUri(this.TrueUri) {
			return
		}
		if global.JsonResEncrypt() {
			this.JsonEncrypt = true
		}
		if global.IsDev() {
			return
		}
		timestamp := this.Ctx.Input.Header("TIMESTAMP")
		random := this.Ctx.Input.Header("RANDOM")
		sign := this.Ctx.Input.Header("SIGN")
		err := authBusiness.CheckSign(timestamp, random, sign)
		if err != nil {
			this.FailReturn(err)
			return
		}
	} else {
		if strings.Index(this.TrueUri, "/internal") == 0 && appId != "20000" {
			this.FailReturn(global.NewError(4004))
			return
		}
		err := authBusiness.CheckAppIdSign(appId, this.Ctx)
		if err != nil {
			this.FailReturn(err)
			return
		}
		//赋予权限
		this.UserId = 1
		this.DyLevel = 3
		this.XhsLevel = 3
		this.TbLevel = 3
	}
}

//校验token
func (this *ApiBaseController) CheckToken() {
	if this.UserId > 0 {
		return
	}
	authBusiness := business.NewAccountAuthBusiness()
	if authBusiness.AuthLoginWhiteUri(this.TrueUri) {
		return
	}
	this.FailReturn(global.NewError(4001))
	return
}

//初始化api
func (this *ApiBaseController) InitApi() {
	this.CheckIp()
	gds := this.Input()
	if len(gds) > 0 {
		this.GetDatas = make(map[string]interface{})
		for k, v := range gds {
			if len(v) > 0 {
				this.GetDatas[k] = v[0]
			}
		}
	}
	if this.Ctx.Request.Method == "POST" || this.Ctx.Request.Method == "PUT" {
		var inputFormat interface{}
		_ = json.Unmarshal(this.Ctx.Input.RequestBody, &inputFormat)
		this.ApiDatas = inputFormat
		if this.ApiDatas == nil {
			this.ApiDatas = make(map[string]interface{}, 0)
		}

		//把get请求的参数也加入到ApiDatas,优先取post
		if apiDatas, ok := this.ApiDatas.(map[string]interface{}); ok {
			for k, v := range this.GetDatas {
				if apiDatas[k] == nil {
					apiDatas[k] = v
				}
			}
			this.ApiDatas = apiDatas
		}

	} else {
		this.ApiDatas = this.GetDatas
	}
	authBusiness := business.NewAccountAuthBusiness()
	this.TrueUri = authBusiness.GetTrueRequestUri(this.Ctx.Input.URI(), this.Ctx.Input.Params())
	this.LogInputOutput("Input", this.ApiDatas)
}

//输入格式化
func (this *ApiBaseController) InputFormat() (retInput global.InputMap) {
	if _, ok := this.ApiDatas.(map[string]interface{}); ok {
		retInput = this.ApiDatas.(map[string]interface{})
	}
	return
}

//输入数据格式化
func (this *ApiBaseController) InputFormatArr() (retInput []global.InputMap) {
	if _, ok := this.ApiDatas.([]interface{}); ok {
		for _, v := range this.ApiDatas.([]interface{}) {
			if _, ok := v.(map[string]interface{}); ok {
				formatV := v.(map[string]interface{})
				retInput = append(retInput, formatV)
			}
		}
	}
	return
}

//记录请求的输入输出
func (this *ApiBaseController) LogInputOutput(logtype string, args interface{}) {
	if global.Cfg.String("request_input_output_log") == "ON" {
		if logtype == "Output" {
			aliLog.LogInput(this.Ctx.Input.Header("X-Request-Id"), this.Ctx.Input.Header("X-Client-Id"), "Output", this.AppId, this.UserId, this.TrueUri, this.Ctx.Request.URL.String(), this.Ctx.Request.Method, this.Ctx.Input.IP(), this.Ctx.Request.UserAgent(), "", args, this.DyLevel, this.XhsLevel, this.TbLevel, this.Ctx.Input.Header("X-Remote-Addr"))
		} else {
			aliLog.LogInput(this.Ctx.Input.Header("X-Request-Id"), this.Ctx.Input.Header("X-Client-Id"), logtype, this.AppId, this.UserId, this.TrueUri, this.Ctx.Request.URL.String(), this.Ctx.Request.Method, this.Ctx.Input.IP(), this.Ctx.Request.UserAgent(), this.Ctx.Request.Referer(), args, this.DyLevel, this.XhsLevel, this.TbLevel, this.Ctx.Input.Header("X-Remote-Addr"))
		}
	}
}

// 处理区间数据
func (c *ApiBaseController) SetIntervalData(esQuery *elasticsearch.ElasticQuery, fieldName string, val string, delta ...int) error {
	condition, err := utils.BuildIntervalCondition(val, delta...)
	if err != nil {
		return err
	}
	esQuery.SetRange(fieldName, condition)
	return nil
}

// 解析区间数据
func (c *ApiBaseController) explainNumberInterval(str string, delta int) (firstNum int, secondNum int, err error) {
	return c.explainNumberInterval(str, delta)
}

func (c *ApiBaseController) GetPage(key ...string) int {
	keyName := "page"
	if len(key) >= 1 {
		keyName = key[0]
	}
	val, err := c.GetInt(keyName, 1)
	if err != nil {
		val = 1
	}
	return val
}

func (c *ApiBaseController) GetPageSize(key string, defSize int, maxSize int) (size int) {
	//var err error
	size = defSize
	data := c.InputFormat()
	size = data.GetInt(key, defSize)
	if size > maxSize {
		size = maxSize
	}
	return
}

func (c *ApiBaseController) GetStringWithQualified(key string, def string, qualified ...string) string {
	str := c.GetString(key, def)
	if str != def {
		if !utils.InArray(str, qualified) {
			c.FailReturn(global.NewError(40000))
			return def
		}
	}
	return str
}

func (c *ApiBaseController) HandleError(err error, code ...int) {
	if err == nil {
		return
	}
	logs.Error("handle err:", err)
	errCode := 50000
	if len(code) > 0 {
		errCode = code[0]
	}
	c.FailReturn(global.NewError(errCode))
}

// 表单校验
func (this *ApiBaseController) FormVerify(input interface{}) {
	valid := validation.Validation{}
	b, _ := valid.Valid(input)
	if !b {
		arr := strings.Split(valid.Errors[0].Key, ".") //切分例如Password.MinSize
		st := reflect.TypeOf(input).Elem()             // 反射获取 input 信息
		field, _ := st.FieldByName(arr[0])             // 获取 Password 参数信息
		this.FailReturn(global.NewMsgError(field.Tag.Get("chn") + valid.Errors[0].Message))
		return
	}
}

//检测抖音用户权限
func (this *ApiBaseController) CheckDyUserGroupRight(minAuthShow, maxAuthShow int) {
	this.MaxTotal = minAuthShow
	if this.DyLevel > 0 {
		this.MaxTotal = maxAuthShow
		this.HasAuth = true
	}
	return
}

//获取当前请求  最低需要的权限
func (this *ApiBaseController) GetMinLevel() int {
	return 0
}

//func (this *ApiBaseController) IsSEOSpider() bool {
//	// 因为seo spider请求会经前端SEO程序转发，直接判断入网IP即可
//	if this.Ctx.Request.RemoteAddr == "47.103.153.227" || this.Ctx.Input.IP() == "47.103.153.227" {
//		return true
//	}
//	return false
//}

// 验证拦截
func (this *ApiBaseController) AsfCheck() {
	disabled := global.Cache.Get(cache.GetCacheKey(cache.SecurityVerifyDisabled))
	// 如果关闭，直接返回
	if disabled == "1" {
		return
	}
	if this.IsMonitor {
		return
	}
	verifyUser := ""
	if this.UserId > 0 {
		verifyUser = global.Cache.Get(cache.GetCacheKey(cache.SecurityVerifyCodeUid, this.UserId))
	} else {
		verifyIp := global.Cache.Get(cache.GetCacheKey(cache.SecurityVerifyCodeIp, this.Ip))
		if verifyIp == "verify" || verifyUser == "verify" {
			//if this.Ip == "47.103.153.227" {
			//	return
			//}
			this.FailReturn(global.NewError(80000))
			return
		}
	}
}

func (this *ApiBaseController) GetUri() string {
	r := this.Ctx.Request.RequestURI
	r, _ = url.QueryUnescape(r)
	//r = strings.Replace(r, "%2f", "/", -1)
	reg := regexp.MustCompile(`(\/)+`)
	r = fmt.Sprintf("%s", reg.ReplaceAllString(r, "/"))
	i := strings.Index(r, "?")
	if i > 0 {
		r = r[:i]
	}
	r = strings.TrimRight(r, "/")
	return r
}

func (this *ApiBaseController) ExportCsv(export, filename string) {
	this.Ctx.Output.Header("Content-Type", "application/csv; charset=gdk")
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	this.Ctx.WriteString("\xEF\xBB\xBF" + export)
}

func (this *ApiBaseController) GetAppId() int {
	return this.AppId
}

func (c *ApiBaseController) DownloadBuf(fileName string, buf *bytes.Buffer) error {
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fileName+"; filename*=utf-8''")
	c.Ctx.Output.Header("Content-Description", "File Transfer")
	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Header("Content-Transfer-Encoding", "binary")
	c.Ctx.Output.Header("Expires", "0")
	c.Ctx.Output.Header("Cache-Control", "must-revalidate")
	c.Ctx.Output.Header("Pragma", "public")
	return c.Ctx.Output.Body(buf.Bytes())
}

//注册登录事件
func (c *ApiBaseController) RegisterLogin(token string, expire_time int64) {
	tokenCookie := c.Ctx.GetCookie(global.LOGINCOOKIENAME)
	if tokenCookie == "" || tokenCookie != token {
		domain := global.Cfg.String("oauth2_cookie_domain")
		expire := expire_time - utils.Time()
		c.Ctx.SetCookie(global.LOGINCOOKIENAME, token, expire, "/", domain, true, true)
	}
}

//注册登录事件
func (c *ApiBaseController) RegisterLogout() {
	domain := global.Cfg.String("oauth2_cookie_domain")
	c.Ctx.SetCookie(global.LOGINCOOKIENAME, "", -1, "/", domain, true, true)
}

func (c *ApiBaseController) Echo(body string, code ...int) {
	w := c.Ctx.ResponseWriter
	statusCode := http.StatusOK
	if len(code) > 0 {
		statusCode = code[0]
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(body))
	c.StopRun()
}

func (c *ApiBaseController) CheckCaptcha() (ok bool, err error) {
	//inputData := c.InputFormat()
	//ticket := inputData.GetString("captcha_ticket", "")
	//randStr := inputData.GetString("captcha_randstr", "")
	//captcha := tencent.NewDefaultCaptcha()
	//resp, err := captcha.DescribeCaptcha(ticket, randStr, c.Ip)
	//if err != nil {
	//	return
	//}
	//if resp.Response.CaptchaCode != nil && *resp.Response.CaptchaCode == tencent.CaptchaCodeSuccess {
	//	ok = true
	//}
	return true, nil
}

//获取查询时间并校验
func (receiver *ApiBaseController) GetRangeDate() (startTime, endTime time.Time, commonError global.CommonError) {
	startDay := receiver.Ctx.Input.Param(":start")
	endDay := receiver.Ctx.Input.Param(":end")
	if startDay == "" {
		commonError = global.NewError(4000)
		return
	}
	if endDay == "" {
		endDay = time.Now().Format("2006-01-02")
	}
	pslTime := "2006-01-02"
	startTime, err := time.ParseInLocation(pslTime, startDay, time.Local)
	if err != nil {
		commonError = global.NewError(4000)
		return
	}
	//时间截止至9.1号
	if startTime.Unix() < 1630425600 {
		startTime = time.Unix(1630425600, 0)
	}
	endTime, err = time.ParseInLocation(pslTime, endDay, time.Local)
	if err != nil {
		commonError = global.NewError(4000)
		return
	}
	if startTime.After(endTime) || endTime.After(startTime.AddDate(0, 0, 90)) || endTime.After(time.Now()) {
		commonError = global.NewError(4000)
		return
	}
	return
}

//禁词
func (receiver *ApiBaseController) KeywordBan(keyword string) {
	if utils.InArrayString(keyword, []string{}) {
		receiver.FailReturn(global.NewError(10000))
		return
	}
	return
}
