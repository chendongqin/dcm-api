package dcm

import (
	"time"
)

type DcUserInvite struct {
	Id          int       `xorm:"not null pk autoincr INT(11)"`
	UserPhone   string    `xorm:"not null default ''  VARCHAR(11)"`
	InvitePhone string    `xorm:"not null default ''  VARCHAR(11)"`
	Platform    string    `xorm:"not null default ''  VARCHAR(50)"`
	CreateTime  time.Time `xorm:"DATETIME"`
	UpdateTime  time.Time `xorm:"DATETIME"`
}
