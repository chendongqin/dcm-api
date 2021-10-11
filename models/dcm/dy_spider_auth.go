package dcm

import (
	"time"
)

type DySpiderAuth struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	Uid        string    `xorm:"not null default '' VARCHAR(50)"`
	Nickname   string    `xorm:"not null default '' VARCHAR(50)"`
	Cookies    string    `xorm:"not null TEXT"`
	Sessionid  string    `xorm:"not null default '' VARCHAR(32)"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
	UpdateTime time.Time `xorm:"TIMESTAMP"`
}
