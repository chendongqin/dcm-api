package dcm

import (
	"time"
)

type DcUserDyCollectTag struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	UserId     int       `xorm:"not null default 0 comment('用户id') INT(11)"`
	Name       string    `xorm:"not null default '' comment('名称') VARCHAR(20)"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
	UpdateTime time.Time `xorm:"TIMESTAMP"`
	DeleteTime time.Time `xorm:"TIMESTAMP"`
}
