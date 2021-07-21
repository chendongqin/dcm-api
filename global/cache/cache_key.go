package cache

type KeyName string

const (
	//appid密钥
	AppIdSecret      KeyName = "appid:secret:%s"
	UserInfo         KeyName = "userinfo:data:%d"
	UserPrevTimeLock KeyName = "user:prev:lock:%d"
	//短信验证码
	SmsCodeLimitBySome KeyName = "dcm:sms:limit:%s:%s" //短信发送限制
	SmsCodeVerify      KeyName = "dcm:sms:code:%s:%s"
	// 全局关闭验证
	SecurityVerifyDisabled KeyName = "dcm:security:verify:disabled"
	//触发滑块验证
	SecurityVerifyCodeUid KeyName = "dcm:security:verify:code:uid:%d"
	SecurityVerifyCodeIp  KeyName = "dcm:security:verify:code:ip:%s"
	//用户Platform的唯一token
	UserPlatformUniqueToken KeyName = "user:unique:token:%d:%d" //userId, platformId
	//直播间商品分类数据缓存
	LiveRoomProductCount KeyName = "live:room:product:count:%s"
)
