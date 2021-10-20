package dcm

import (
	"time"
)

type DcVipOrder struct {
	Id             int       `xorm:"not null pk autoincr INT(11)"`
	UserId         int       `xorm:"not null comment('用户id') index INT(11)"`
	Username       string    `xorm:"not null default '' comment('用户手机号') CHAR(11)"`
	TradeNo        string    `xorm:"not null comment('交易订单号') unique CHAR(25)"`
	Channel        int       `xorm:"not null default 0 comment('0后台,1web,2app,第一次支付请求，当微信支付时，不可跨平台支付') TINYINT(1)"`
	InterTradeNo   string    `xorm:"not null default '' comment('第三方饭回交易订单号') VARCHAR(100)"`
	IosReceipt     string    `xorm:"not null default '' comment('苹果支付凭证') index VARCHAR(1000)"`
	OrderType      int       `xorm:"not null comment('订单类型1购买会员2会员续费3协同账号购买4协同账号续费5团队续费6赠送7直播监控购买') TINYINT(1)"`
	PayType        string    `xorm:"not null default '' comment('支付方式') VARCHAR(10)"`
	Platform       string    `xorm:"not null comment('douyin：抖音，xiaohongshu：小红书，taobao：淘宝') ENUM('douyin','taobao','xiaohongshu')"`
	Level          int       `xorm:"not null default 0 comment('购买等级') TINYINT(1)"`
	BuyDays        int       `xorm:"not null default 0 comment('购买天数') SMALLINT(5)"`
	Title          string    `xorm:"not null comment('订单描述标题') VARCHAR(100)"`
	Amount         string    `xorm:"not null default 0.00 comment('订单支付金额') DECIMAL(10,2)"`
	TicketAmount   string    `xorm:"not null default 0.00 comment('优惠券金额') DECIMAL(10,2)"`
	TicketId       int       `xorm:"not null default 0 comment('优惠券id') INT(11)"`
	Status         int       `xorm:"not null default 0 comment('订单状态，1有效，2已取消，0未处理') TINYINT(1)"`
	PayStatus      int       `xorm:"not null default 0 comment('支付状态') TINYINT(1)"`
	GoodsInfo      string    `xorm:"not null comment('商品信息') VARCHAR(800)"`
	Remark         string    `xorm:"not null comment('备注') VARCHAR(100)"`
	ExpirationTime time.Time `xorm:"comment('订单支付过期时间') TIMESTAMP"`
	CreateTime     time.Time `xorm:"comment('创建时间') TIMESTAMP"`
	UpdateTime     time.Time `xorm:"comment('更新时间 ') TIMESTAMP"`
	PayTime        time.Time `xorm:"comment('用户支付回调时间') TIMESTAMP"`
	Referrer       string    `xorm:"not null default '' comment('推荐人') VARCHAR(20)"`
	InvoiceId      int       `xorm:"not null default 0 comment('开票ID') INT(11)"`
	AdminId        int       `xorm:"not null default 0 comment('操作人id') INT(11)"`
	AdminNickName  string    `xorm:"not null default '' comment('操作人昵称') VARCHAR(50)"`
	Image          string    `xorm:"not null default '' comment('相关图片') VARCHAR(2000)"`
	TestOrder      int       `xorm:"not null default 0 comment('是否测试订单') TINYINT(1)"`
}

type DcVipInvoiceOrder struct {
	DcVipOrder        `xorm:"extends"`
	DcVipOrderInvoice `xorm:"extends"`
}
