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
	SmsCodeLimitBySome KeyName = "dcm:sms:limit:%s:%s:%s" //短信发送限制
	SmsCodeVerify      KeyName = "dcm:sms:code:%s:%s"
	// 全局关闭验证
	SecurityVerifyDisabled KeyName = "dcm:security:verify:disabled"
	//触发滑块验证
	SecurityVerifyCodeUid KeyName = "dcm:security:verify:code:uid:%d"
	SecurityVerifyCodeIp  KeyName = "dcm:security:verify:code:ip:%s"
	//用户Platform的唯一token
	UserPlatformUniqueToken KeyName = "dcm:user:unique:token:%d:%d" //userId, platformId
	UserActionLockKey       KeyName = "dcm:user:action:lock:%d"     //userId
	//达人
	AuthorLiveMap KeyName = "dcm:author:map:%s"
	//直播间商品分类数据缓存
	LiveRoomProductCount     KeyName = "dcm:live:room:product:count:%s"
	ProductAuthorCount       KeyName = "dcm:product:author:%s:%s:%s"
	ShopAuthorCount          KeyName = "dcm:shop:author:%s:%s:%s"
	LivePromotionsDetailList KeyName = "dcm:live:promotions:detail:%s:%s"
	LiveRoomProductList      KeyName = "dcm:live:product:list:%s"
	//商品达人数据缓存
	ProductAuthorAllList      KeyName = "dcm:product:author:row:%s:%s"
	ProductAuthorAllMap       KeyName = "dcm:product:author:info:%s:%s"
	ProductAwemeAuthorCount   KeyName = "dcm:product:aweme:author:%s:%s:%s"
	ShopAwemeAuthorCount      KeyName = "dcm:shop:aweme:author:%s:%s:%s"
	ProductAwemeAuthorAllList KeyName = "dcm:product:aweme:author:row:%s:%s"
	ProductAwemeAuthorAllMap  KeyName = "dcm:product:aweme:author:info:%s:%s"
	AuthorProductAllList      KeyName = "dcm:author:product:row:%s:%s"
	ProductAuthorAwemesList   KeyName = "dcm:product:author:awemes:%s:%s"
	AuthorViewProductAllList  KeyName = "dcm:author:view:product:%s:%s:%s"
	RedAuthorRooms            KeyName = "dcm:red:author:room:%s"
	RedAuthorLivingRooms      KeyName = "dcm:red:author:living:room"
	RedAuthorMapCache         KeyName = "dcm:red:author:map:%s"
	//小店
	ShopProductAnalysisScanList      KeyName = "dcm:shop:product:analysis:%s:%s:%s"
	ShopProductAnalysisCountScanList KeyName = "dcm:shop:product:analysis:count:%s:%s:%s"
	ShopLiveAuthorAllList            KeyName = "dcm:shop:live:author:row:%s"
	ShopAwemeAuthorAllList           KeyName = "dcm:shop:aweme:author:row:%s"
	//视频商品列表
	AwemeProductByDate KeyName = "dcm:aweme:product:%s:%s:%s"
	//商品视频列表
	ProductRelateAweme KeyName = "dcm:product:aweme:%s:%s:%s:%s"
	//榜单数据缓存
	DyRankCache KeyName = "dcm:rank:%s:%s"
	//爬虫加速限制频次
	SpiderSpeedUpLimit KeyName = "dcm:spider:limit:%s:%s" //spidername,authorId

	//脚本锁
	DyMonitorUpdateRoomLock KeyName = "dcm:cmd:monitor:update:room:%s"

	//定时任务的键
	AmountExpireWechatNotice KeyName = "dcm:cmd:account:expire:notice:%d" //天数

	//腾讯广告
	TencentAdAccessToken  KeyName = "dcm:tencent:ad:token"
	TencentAdRefreshToken KeyName = "dcm:tencent:ad:fresh"

	//爬虫接口达人搜索
	SpiderAuthorSearchKeyWord KeyName = "dcm:spider:author:search:%s"
)
