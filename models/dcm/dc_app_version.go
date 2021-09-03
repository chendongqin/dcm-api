package dcm

import (
	"time"
)

type DcAppVersion struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	Platform   int       `xorm:"not null default 0 comment('1ios,2安卓') TINYINT(1)"`
	Version    string    `xorm:"not null comment('版本') VARCHAR(10)"`
	Info       string    `xorm:"not null comment('版本信息') VARCHAR(500)"`
	Force      int       `xorm:"not null default 0 comment('是否强制更新') TINYINT(1)"`
	Url        string    `xorm:"not null default '' comment('下载链接') VARCHAR(255)"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
	UpdateTime time.Time `xorm:"TIMESTAMP"`
}
