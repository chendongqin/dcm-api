package cache

type KeyName string

const (

	//短信验证码
	SmsCodeLimitBySome KeyName = "dcm:sms:limit:%s:%s" //短信发送限制
	SmsCodeVerify KeyName = "dcm:sms:code:%s:%s"
	// 全局关闭验证
	SecurityVerifyDisabled KeyName = "dcm:security:verify:disabled"
	//触发滑块验证
	SecurityVerifyCodeUid KeyName = "dcm:security:verify:code:uid:%d"
	SecurityVerifyCodeIp  KeyName = "dcm:security:verify:code:ip:%s"
	//用户Platform的唯一token
	UserPlatformUniqueToken KeyName = "user:unique:token:%d:%d" //userId, platformId

)
