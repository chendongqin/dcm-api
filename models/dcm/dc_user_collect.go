package dcm

import (
	"time"
)

type DcUserCollect struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	UserId     int       `xorm:"not null default 0 comment('用户id') index(UN_INDEX) INT(11)"`
	Platform   int       `xorm:"not null default 1 comment('1抖音2小红书3淘宝') index(UN_INDEX) TINYINT(1)"`
	AuthorId   string    `xorm:"not null default '' comment('达人id') index(UN_INDEX) VARCHAR(50)"`
	Label      string    `xorm:"not null default '' comment('达人分类') VARCHAR(20)"`
	TagId      int       `xorm:"not null default 0 comment('分组id') INT(11)"`
	Remark     string    `xorm:"not null default '' comment('备注') VARCHAR(100)"`
	Status     int       `xorm:"not null default 1 comment('状态0删除1正常') TINYINT(1)"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
	UpdateTime time.Time `xorm:"TIMESTAMP"`
}
