package dcm

import (
	"time"
)

type DcUserVip struct {
	Id             int       `xorm:"not null pk autoincr INT(11)"`
	UserId         int       `xorm:"not null default 0 unique(USER_LEVEL) INT(11)"`
	Platform       int       `xorm:"not null default 1 comment('1抖音2小红书3淘宝') unique(USER_LEVEL) TINYINT(1)"`
	Level          int       `xorm:"not null default 0 comment('等级：0普通，1vip，2svip,3专业版') TINYINT(1)"`
	Expiration     time.Time `xorm:"comment('过期时间') TIMESTAMP"`
	OrderLevel     int       `xorm:"not null default 0 comment('订单暂缓有效等级') TINYINT(1)"`
	OrderValidDay  int       `xorm:"not null default 0 comment('订单暂缓有效天数') SMALLINT(5)"`
	UpdateTime     time.Time `xorm:"comment('更新时间') TIMESTAMP"`
	ValueType      int       `xorm:"not null default 0 comment('0无效1赠送2购买3子账号') TINYINT(1)"`
	ParentId       int       `xorm:"not null default 0 comment('主账户id') INT(11)"`
	SubIds         string    `xorm:"not null default '' comment('子账号集合') VARCHAR(255)"`
	SubNum         int       `xorm:"not null default 0 comment('子账号数') SMALLINT(3)"`
	SubExpiration  time.Time `xorm:"comment('子账号过期时间') TIMESTAMP"`
	LiveMonitorNum int       `xorm:"not null default 0 comment('购买的直播监控次数，长期有效') SMALLINT(5)"`
	Remark        string    `xorm:"not null default '' comment('备注') VARCHAR(30)"`
}
