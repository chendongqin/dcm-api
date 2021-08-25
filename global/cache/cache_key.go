package cache

type KeyName string

const (
	//appid密钥
	AppIdSecret KeyName = "dcm:appid:secret:%s"
	//配置缓存
	ConfigKeyCache         KeyName = "dcm:config:cache:%s"
	LongTimeConfigKeyCache KeyName = "dcm:config:long:cache"
	//用户相关
	UserInfo         KeyName = "dcm:user:info:data:%d"
	UserLevel        KeyName = "dcm:user:level:%d:%d"
	UserPrevTimeLock KeyName = "dcm:user:prev:lock:%d"
	UserActionLock   KeyName = "dcm:user:action:lock:%s:%s"
	//短信验证码
	SmsCodeLimitBySome KeyName = "dcm:sms:limit:%s:%s" //短信发送限制
	SmsCodeVerify      KeyName = "dcm:sms:code:%s:%s"
	// 全局关闭验证
	SecurityVerifyDisabled KeyName = "dcm:security:verify:disabled"
	//触发滑块验证
	SecurityVerifyCodeUid KeyName = "dcm:security:verify:code:uid:%d"
	SecurityVerifyCodeIp  KeyName = "dcm:security:verify:code:ip:%s"
	//用户Platform的唯一token
	UserPlatformUniqueToken KeyName = "dcm:user:unique:token:%d:%d" //userId, platformId
	UserActionLockKey       KeyName = "dcm:user:action:lock:%d"     //userId
	//直播间商品分类数据缓存
	LiveRoomProductCount KeyName = "dcm:live:room:product:count:%s"
	ProductAuthorCount   KeyName = "dcm:product:author:%s:%s:%s"
	//商品达人数据缓存
	ProductAuthorAllList     KeyName = "dcm:product:author:row:%s:%s"
	AuthorProductAllList     KeyName = "dcm:author:product:row:%s:%s"
	AuthorViewProductAllList KeyName = "dcm:author:view:product:%s:%s"
	//榜单数据缓存
	DyRankCache              KeyName = "dcm:rank:%s:%s"
	DyAuthorTakeGoodsRank    KeyName = "dcm:rank:author:goods:%s:%s:%s:%s:%s:%s:%s:%s"
	DyAuthorFollowerRank     KeyName = "dcm:rank:author:follower:%s:%s:%s"
	DyAwemeShareRank         KeyName = "dcm:rank:awemer:share:%s"
	DyLiveShareWeekRank      KeyName = "dcm:rank:live:week:share:%s"
	DyLiveHourPopularityRank KeyName = "dcm:rank:live:hour:popularity:%s"
	DyLiveHourSellRank       KeyName = "dcm:rank:live:hour:sell:%s"
	DyLiveTopRank            KeyName = "dcm:rank:live:top:%s"
	DyLiveHourRank           KeyName = "dcm:rank:live:hour%s"
	DyStartAuthorLiveRank    KeyName = "dcm:rank:star:author:live:%s"
	DyStartAuthorVideoRank   KeyName = "dcm:rank:star:author:video:%s"
)
