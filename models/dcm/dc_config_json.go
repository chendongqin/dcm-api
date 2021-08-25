package dcm

import (
	"time"
)

type DcConfigJson struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	KeyName    string    `xorm:"not null default '' comment('键名唯一') unique VARCHAR(30)"`
	Value      string    `xorm:"comment('值') TEXT"`
	Note       string    `xorm:"not null default '' comment('备注') VARCHAR(100)"`
	CreateTime time.Time `xorm:"DATETIME"`
	UpdateTime time.Time `xorm:"DATETIME"`
	Auth       int       `xorm:"not null default 0 comment('0对外1需要登录2内部配置') TINYINT(255)"`
}
