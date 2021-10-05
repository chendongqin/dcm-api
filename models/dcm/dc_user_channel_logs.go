package dcm

import (
	"time"
)

type DcUserChannelLogs struct {
	Id          int       `xorm:"not null pk autoincr INT(11)"`
	UserId      int       `xorm:"not null comment('用户id') INT(11)"`
	Channel     string    `xorm:"not null default '' comment('渠道') index(CHNANEL_INDEX) VARCHAR(100)"`
	ChannelWord string    `xorm:"not null default '' comment('渠道关键词') index(CHNANEL_INDEX) VARCHAR(100)"`
	AppId       int       `xorm:"not null default 0 comment('10000：pc端，10001：h5，10002:微信小程序,10003、10004：app,10005:Wap') INT(11)"`
	Ip          string    `xorm:"not null default '' comment('ip') VARCHAR(50)"`
	CreateTime  time.Time `xorm:"comment('点击创建时间') TIMESTAMP"`
}
