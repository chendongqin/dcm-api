package dcm

import (
	"time"
)

type DcLiveMonitorRoom struct {
	Id           int       `xorm:"not null pk autoincr INT(11)"`
	MonitorId    int       `xorm:"not null default 0 comment('关联监测记录ID') index index(USER_ID) INT(11)"`
	UserId       int       `xorm:"not null comment('冗余用户ID') index(USER_ID) INT(11)"`
	AuthorId     string    `xorm:"not null default '' comment('冗余达人ID') VARCHAR(64)"`
	RoomId       string    `xorm:"not null default '' comment('直播间ID') VARCHAR(64)"`
	OpenId       string    `xorm:"not null default '' comment('微信ID') VARCHAR(64)"`
	Status       int       `xorm:"not null default 0 comment('状态 2 在播 4 下播 5 异常') TINYINT(1)"`
	FinishTime   int       `xorm:"not null default 0 comment('下播时间戳') INT(11)"`
	FinishNotice int       `xorm:"not null default 0 comment('下播提醒') index(CREATED_TIME) TINYINT(1)"`
	ProductId    string    `xorm:"not null default '' comment('商品ID，逗号分割') VARCHAR(512)"`
	CreateTime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') index(CREATED_TIME) TIMESTAMP"`
	UpdateTime   time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('更新时间') TIMESTAMP"`
	Gmv          string    `xorm:"not null default 0.00 comment('销售额') DECIMAL(10,2)"`
	UserTotal    int       `xorm:"not null default 0 comment('观看人次') INT(11)"`
	Sales        int       `xorm:"not null default 0 comment('销量') INT(11)"`
}
