package controllers

import (
	"bytes"
	"dongchamao/controllers"
	"dongchamao/global"
	"dongchamao/global/cache"
	"dongchamao/global/utils"
	"dongchamao/models/business"
	"dongchamao/models/dcm"
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
)

type ApiBaseController struct {
	GetDatas         map[string]interface{}
	ApiDatas         interface{}
	UserId           int
	UserInfo         dcm.DcUser
	IsMonitor        bool
	appId            int
	appSecret        string
	Ip               string
	IsInitToken      bool               //是否初始化过token
	LastInitTokenErr global.CommonError //记录首次初始化token的错误
	Token            string
	controllers.BaseController
}

type AppData struct {
	appId     int
	appSecret string
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
	//this.Ip = this.Ctx.Input.IP()
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

//校验签名
func (this *ApiBaseController) CheckSign() {
	InputDatas := this.InputFormat()
	this.appId = InputDatas.GetInt("appId", 0)
	if this.appId == 10000 || this.appId == 10001 || this.appId == 10003 || this.appId == 10004 {
		return
	}

	//签名过期
	currentTime := utils.Time()
	timeStamp := int64(InputDatas.GetInt("timeStamp", 0))

	logs.Info(InputDatas)
	logs.Info(currentTime)
	logs.Info(timeStamp)

	if currentTime > timeStamp+1800 || currentTime < timeStamp-600 {
		this.FailReturn(global.NewError(40002))
		return
	}
	//签名校验
	//appData := logic.NewCacheHandle().GetApp(this.appId)
	////if appData,ok := appDatas[this.appId];ok{
	//if appData != nil {
	//	signMap := make(map[string]string)
	//	for k, v := range InputDatas {
	//		if k != "sign" {
	//			signMap[k] = utils.InterfaceToString(v)
	//		}
	//	}
	//	sign := InputDatas.GetString("sign", "")
	//	realSign := utils.Wechat_makeSign(signMap, appData.Appsecret)
	//	if realSign != sign {
	//		this.FailReturn(global.NewError(40003))
	//		return
	//	}
	//} else {
	//	this.FailReturn(global.NewError(40001))
	//	return
	//}
}

func (this *ApiBaseController) CheckAppAccess() {
	this.appId = utils.ParseInt(strings.Trim(this.Ctx.Input.Header("AccessId"), ""), 0)
	this.appSecret = strings.Trim(this.Ctx.Input.Header("AccessKey"), "")
	//appData := logic.NewCacheHandle().GetApp(this.appId)
	//if appData == nil || this.appSecret != appData.Appsecret {
	//	this.FailReturn(global.NewError(40001))
	//	return
	//}
}

//模拟真实访问构造的token
func (this *ApiBaseController) checkMonitorToken() {
	//apiMonitor := monitor.NewApiMonitor()
	//checkSign := apiMonitor.BuildSign(this.Ctx.Request.URL.Path)
	//if checkSign != this.Ctx.Input.Header("X-Monitor-Sign") {
	//	this.FailReturn(global.NewError(40004))
	//	return
	//}
	//this.UserId = 88905
	//this.IsMonitor = true
	//uinfo := apiv1models.NewUserModel()
	//err := uinfo.GetInfo(this.UserId, true)
	//if err != nil || uinfo.Id == 0 {
	//	this.FailReturn(global.NewError(42004))
	//	return
	//}
	//this.UserInfo = uinfo
	//this.UserInfo.GroupId = 4
}

//不需要手机号码账号就能访问的接口白名单
var NoPhoneWhiteRoute = []string{
	"/v1/discountActivity/coupon/couponAddScore",
	"/v1/vip/order/createAppleMonitorOrder",
	"/v1/vip/order/createAppleOrder",
	"/v1/vip/order/getOrderPrice",
}

func (this *ApiBaseController) InitUserToken(args ...bool) (commonErr global.CommonError) {
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

	//如果是监控发过来的请求，则通过checkMonitorToken进行校验
	if this.Ctx.Input.Header("X-Monitor-Sign") != "" {
		this.checkMonitorToken()
		return nil
	}
	//this.UserInfo = &apiv1models.VoUser{
	//	GroupId: 0,
	//}

	tokenString := this.Ctx.Input.Cookie(global.LOGINCOOKIENAME)
	//cookie没有身份信息  从头部获取
	if tokenString == "" {
		tokenString = strings.Replace(this.Ctx.Input.Header("Authorization"), "Bearer ", "", 1)
	}
	if tokenString == "" {
		return global.NewError(40004)
	}
	token, _ := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Cfg.String("auth_code")), nil
	})
	if token == nil {
		return global.NewError(40015)
	}
	userBusiness := business.NewUserBusiness()
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		json.Marshal(claims)
		jsonstr, err0 := json.Marshal(claims)
		tokenStruct := &business.TokenData{}
		err0 = json.Unmarshal(jsonstr, tokenStruct)
		if err0 != nil {
			return global.NewError(40004)
		}
		expire_time := tokenStruct.ExpireTime
		if utils.Time() > expire_time {
			return global.NewError(40005)
		} else {
			this.UserId = tokenStruct.Id
		}
		if this.UserId == 0 {
			return global.NewError(42004)
		}
		this.appId = tokenStruct.AppId
		this.Token = tokenString
		userInfo := dcm.DcUser{}
		_,err := dcm.Get(this.UserId, &userInfo)
		if err != nil || userInfo.Id == 0 {
			return global.NewError(42004)
		}
		this.UserInfo = userInfo
		//判断用户状态
		if this.UserInfo.Status == 0 {
			return global.NewError(45000)
		}
		//除bindphone外的接口 没有phone不让访问
		if len(args) <= 0 || (len(args) == 1 && args[0] == true) {
			if utils.InArray(this.GetUri(), NoPhoneWhiteRoute) == false {
				if this.UserInfo.Username == "" {
					return global.NewError(42014)
				}
			}
		}

		//处理连续登录次数统计处理
		//if this.UserInfo.UpdateVisitedTimes() {
		//	//发生更新，才需要清除缓存
		//	logic.NewCacheHandle().DeleteUserInfoCache(this.UserInfo.Id)
		//}

		//验证 user Platform token的唯一性
		uniqueToken, exist := userBusiness.GetUniqueToken(int(this.UserId), this.appId, true)
		if exist == false {
			return global.NewError(40005)
		} else {
			if tokenString != uniqueToken {
				//cmmlog.LoginLog(this.Ctx.Input.Header("X-Request-Id"), this.Ctx.Input.Header("X-Client-Id"), this.appId, this.UserId, "current:"+tokenString+"|unique:"+uniqueToken, this.Ctx.Input.IP(), this.Ctx.Request.UserAgent(), "unique", "unique_loginout")

				this.RegisterSignOut()
				return global.NewError(40006)
			}
		}

		if this.appId == 10000 || this.appId == 10001 {
			this.RegisterSignin(tokenString, expire_time)
		}

		//异步记录用户行为日志
		this.LogInputOutput("Format", this.ApiDatas)

	} else {
		return global.NewError(40004)
	}
	return
}

//校验token
func (this *ApiBaseController) CheckToken(args ...bool) {
	err := this.InitUserToken(args...)
	if err != nil {
		this.FailReturn(err)
		return
	}
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
	if this.Ctx.Request.Method == "POST" || this.Ctx.Request.Method == "PUT"{
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
	//if global.Cfg.String("request_input_output_log") == "ON" {
	//	userGroupId := this.GetUserGroupId()
		//if logtype == "Output" {
		//	cmmlog.LogInput(this.Ctx.Input.Header("X-Request-Id"), this.Ctx.Input.Header("X-Client-Id"), "Output", this.appId, this.UserId, this.Ctx.Request.URL.String(), this.Ctx.Request.Method, this.Ctx.Input.IP(), this.Ctx.Request.UserAgent(), "", args, userGroupId, this.Ctx.Input.Header("X-Remote-Addr"))
		//} else {
		//	cmmlog.LogInput(this.Ctx.Input.Header("X-Request-Id"), this.Ctx.Input.Header("X-Client-Id"), logtype, this.appId, this.UserId, this.Ctx.Request.URL.String(), this.Ctx.Request.Method, this.Ctx.Input.IP(), this.Ctx.Request.UserAgent(), this.Ctx.Request.Referer(), args, userGroupId, this.Ctx.Input.Header("X-Remote-Addr"))
		//}
	//}
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

//func (c *ApiBaseController) UserRightLimit(query *logic.PageQuery, maxLength int) bool {
//	if maxLength == 0 {
//		return false
//	}
//	//最大可查看条数限制
//	maxPage := int(math.Ceil(float64(maxLength) / float64(query.Size)))
//	if query.Page > maxPage {
//		return false
//	}
//	if query.Page == maxPage && (query.Size*maxPage) >= maxLength {
//		query.SetSize(maxLength - query.Size*(maxPage-1))
//	}
//	return true
//}

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

////beego validate 设置成中文
//func (c *ApiBaseController) setVerifyMessage() {
//	// 设置表单验证messages
//	var MessageTmpls = map[string]string{
//		"Required":     "不能为空",
//		"Min":          "最小值 为 %d",
//		"Max":          "最大值 为 %d",
//		"Range":        "范围 为 %d 到 %d",
//		"MinSize":      "最短长度 为 %d",
//		"MaxSize":      "最大长度 为 %d",
//		"Length":       "长度必须 为 %d",
//		"Alpha":        "必须是有效的字母",
//		"Numeric":      "必须是有效的数字",
//		"AlphaNumeric": "必须是有效的字母或数字",
//		"Match":        "必须匹配 %s",
//		"NoMatch":      "必须不匹配 %s",
//		"AlphaDash":    "必须是有效的字母、数字或连接符号(-_)",
//		"Email":        "必须是有效的电子邮件地址",
//		"IP":           "必须是有效的IP地址",
//		"Base64":       "必须是有效的base64字符",
//		"Mobile":       "必须是有效的手机号码",
//		"Tel":          "必须是有效的电话号码",
//		"Phone":        "必须是有效的电话或移动电话号码",
//		"ZipCode":      "必须是有效的邮政编码",
//	}
//
//	validation.SetDefaultMessage(MessageTmpls)
//}

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

//检测用户权限
func (this *ApiBaseController) CheckUserGroupRight() {

}

//获取当前请求  最低需要的权限
func (this *ApiBaseController) GetMinGroupId() int {
	return 0
}

func (this *ApiBaseController) Prepare() {
	this.InitApi()
	this.CheckUserGroupRight()
	this.AsfCheck()
}

func (this *ApiBaseController) IsSEOSpider() bool {
	// 因为seo spider请求会经前端SEO程序转发，直接判断入网IP即可
	if this.Ctx.Request.RemoteAddr == "47.103.153.227" || this.Ctx.Input.IP() == "47.103.153.227" {
		return true
	}
	return false
}

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
			if this.Ip == "47.103.153.227" {
				return
			}
			this.FailReturn(global.NewError(80000))
			return
		}
	}
}

func (this *ApiBaseController) GetUserGroupId() int {
	groupId := 0
	//if this.UserInfo != nil && this.UserInfo.Status == logic.UserStatusNormal {
	//	groupId = this.UserInfo.GroupId
	//}
	return groupId
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
	return this.appId
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
func (c *ApiBaseController) RegisterSignin(token string, expire_time int64) {
	tokenCookie := c.Ctx.GetCookie(global.LOGINCOOKIENAME)
	if tokenCookie == "" || tokenCookie != token {
		domain := global.Cfg.String("oauth2_cookie_domain")
		expire := expire_time - utils.Time()
		c.Ctx.SetCookie(global.LOGINCOOKIENAME, token, expire, "/", domain, true, true)
	}
}

//注册登录事件
func (c *ApiBaseController) RegisterSignOut() {
	domain := global.Cfg.String("oauth2_cookie_domain")
	c.Ctx.SetCookie(global.LOGINCOOKIENAME, "", -1, "/", domain, true, true)
	//c.Ctx.SetCookie(global.LOGINCOOKIENAME, "", -1,"/", domain,false,false)
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
	return true,nil
}
