package dcm

import (
	"time"
)

type DcUserToken struct {
	Id          int       `xorm:"not null pk autoincr INT(11)"`
	UserId      int       `xorm:"not null default 0 comment('用户id') unique(userIdPlatform) INT(11)"`
	AppPlatform int       `xorm:"not null default 0 comment('1:pc 2:h5 :3:小程序 4:app') unique(userIdPlatform) INT(11)"`
	Token       string    `xorm:"not null default '' comment('当前唯一可用的token') index VARCHAR(500)"`
	CreateTime  time.Time `xorm:"TIMESTAMP"`
	UpdateTime  time.Time `xorm:"TIMESTAMP"`
}
