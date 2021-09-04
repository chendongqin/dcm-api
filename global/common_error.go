package global

type CommonError interface {
	Error() (int, string)
}

type CommonErrors struct {
	errCode int
	errMsg  string
}

var ErrCode = map[int]string{
	4000: "参数错误",
	//登陆相关
	4001: "签名无效",
	4002: "签名已过期",
	4003: "您的账号在另一个设备登录，若非本人操作，请及时更改密码",
	4004: "您没有该接口的权限",
	4005: "您需要绑定手机号才可继续操作",
	4006: "您未扫描二维码",
	//用户相关
	4200: "账号或密码不能为空",
	4201: "账号格式错误",
	4202: "账号已存在,无法注册",
	4203: "创建账号失败",
	4204: "账号不存在",
	4205: "账号格式错误",
	4206: "会员类型错误",
	4207: "两次密码不一致",
	4208: "账号或密码错误",
	4209: "验证码错误",
	4210: "密码长度应在6-24位范围内",
	4211: "操作过于频繁",
	4212: "账号已被禁用",
	4213: "账号修改密码失败",
	4214: "旧密码错误",
	4215: "账号已存在,无法绑定",

	//微信相关
	4300: "未关注官方微信",
	4301: "unionid不能为空",
	4302: "微信登录失败",
	4303: "unionid不能为空",
	4304: "不存在该微信信息",
	4305: "该手机已经绑定微信",

	4040: "资源不存在",
	4041: "签名错误",

	5000: "系统错误",
	6000: "操作过于频繁,请稍后重试",
	// 爬虫相关
	6100: "更新成功，请在30分钟后刷新页面查看最新数据",

	8000:  "滑块验证",
	10000: "由于相关法律和政策，无法展示相关结果",
	88888: "产品上线倒计时",
}

func NewError(errCode int) CommonError {
	if errCode == 0 {
		return nil
	}
	errMsg := ErrCode[errCode]
	return &CommonErrors{
		errCode: errCode,
		errMsg:  errMsg,
	}
}

func NewCodeError(errCode int, errMsg string) CommonError {
	return &CommonErrors{
		errCode: errCode,
		errMsg:  errMsg,
	}
}

func NewCommonError(err error) CommonError {
	errMsg := err.Error()
	return &CommonErrors{
		errCode: 5000,
		errMsg:  errMsg,
	}
}

func NewMsgError(errMsg string) CommonError {
	return &CommonErrors{
		errCode: 5000,
		errMsg:  errMsg,
	}
}

func (this *CommonErrors) Error() (int, string) {
	return this.errCode, this.errMsg
}

type Error struct {
	Msg string
}

func (e *Error) Error() string {
	return e.Msg
}

func NewNormalError(msg string) *Error {
	return &Error{Msg: msg}
}

type TimeoutError struct {
}

func (e *TimeoutError) Error() string {
	return "read cache timeout"
}
