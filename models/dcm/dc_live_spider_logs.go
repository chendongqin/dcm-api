package dcm

import (
	"time"
)

type DcLiveSpiderLogs struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	AuthorId   string    `xorm:"not null default '' VARCHAR(64)"`
	Top        int       `xorm:"not null default 0 SMALLINT(5)"`
	AddLog     string    `xorm:"not null default '' VARCHAR(20)"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
}
