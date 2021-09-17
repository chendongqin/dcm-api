package dcm

import (
	"time"
)

type DcUserKeywordsRecord struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	Keyword    string    `xorm:"not null default '' comment('搜索关键词') unique VARCHAR(50)"`
	Count      int       `xorm:"not null default 0 comment('搜索次数') INT(11)"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
	UpdateTime time.Time `xorm:"TIMESTAMP"`
}
