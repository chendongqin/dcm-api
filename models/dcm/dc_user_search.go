package dcm

import (
	"time"
)

type DcUserSearch struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	UserId     int       `xorm:"not null default 0 comment('用户ID') index(USER_INDEX) INT(11)"`
	SearchType string    `xorm:"not null comment('筛选查询类型') index(USER_INDEX) ENUM('author','live','mcn','product','shop')"`
	Note       string    `xorm:"not null default '' comment('备注') VARCHAR(50)"`
	Content    string    `xorm:"not null default '' comment('json内容串') VARCHAR(2000)"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
	UpdateTime time.Time `xorm:"TIMESTAMP"`
}
