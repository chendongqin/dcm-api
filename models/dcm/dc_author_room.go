package dcm

import (
	"time"
)

type DcAuthorRoom struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	AuthorId   string    `xorm:"not null comment('达人id') VARCHAR(64)"`
	Nickname   string    `xorm:"not null default '' comment('达人昵称') VARCHAR(50)"`
	UniqueId   string    `xorm:"not null default '' comment('抖音号') VARCHAR(50)"`
	LivingTime time.Time `xorm:"comment('直播预告时间') TIMESTAMP"`
	Status     int       `xorm:"not null default 1 comment('状态0不上架1正常') TINYINT(1)"`
	CreateTime time.Time `xorm:"comment('创建时间') TIMESTAMP"`
	UpdateTime time.Time `xorm:"comment('更新时间') TIMESTAMP"`
	RoomId     string    `xorm:"not null default '' comment('直播间id') VARCHAR(64)"`
	Weight     int       `xorm:"not null default 0 comment('权重') SMALLINT(5)"`
}
