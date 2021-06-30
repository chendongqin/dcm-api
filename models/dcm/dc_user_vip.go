package dcm

import (
	"time"
)

type DcUserVip struct {
	Id            int       `xorm:"not null pk autoincr INT(11)"`
	UserId        int       `xorm:"not null default 0 INT(11)"`
	Platform      int       `xorm:"not null default 1 comment('1抖音2小红书3淘宝') TINYINT(1)"`
	Level         int       `xorm:"not null default 0 comment('等级：0普通，1vip，2svip,3专业版') TINYINT(1)"`
	Expiration    time.Time `xorm:"comment('过期时间') TIMESTAMP"`
	OrderLevel    int       `xorm:"not null default 0 comment('订单暂缓有效等级') TINYINT(1)"`
	OrderValidDay int       `xorm:"not null default 0 comment('订单暂缓有效天数') SMALLINT(5)"`
	UpdateTime    time.Time `xorm:"comment('更新时间') TIMESTAMP"`
}
