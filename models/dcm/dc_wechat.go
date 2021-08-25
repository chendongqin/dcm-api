package dcm

import (
	"time"
)

type DcWechat struct {
	Id              int       `xorm:"not null pk autoincr INT(11)"`
	Openid          string    `xorm:"not null default '' VARCHAR(128)"`
	Unionid         string    `xorm:"not null default '' VARCHAR(128)"`
	NickName        string    `xorm:"not null default '' comment('昵称') VARCHAR(128)"`
	Avatar          string    `xorm:"not null default '' comment('头像') VARCHAR(255)"`
	Sex             int       `xorm:"not null default 0 comment('用户的性别，值为1时是男性，值为2时是女性，值为0时是未知') TINYINT(3)"`
	Country         string    `xorm:"not null default '' VARCHAR(128)"`
	Province        string    `xorm:"not null default '' VARCHAR(128)"`
	City            string    `xorm:"not null default '' VARCHAR(128)"`
	Language        string    `xorm:"not null default '' VARCHAR(128)"`
	Remark          string    `xorm:"not null default '' comment('备注') VARCHAR(255)"`
	Subscribe       int       `xorm:"not null default 0 comment('用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。') TINYINT(3)"`
	SubscribeTime   int64     `xorm:"not null default 0 comment('关注时间') INT(11)"`
	UnsubscribeTime int64     `xorm:"not null default 0 comment('取消关注时间') INT(11)"`
	SubscribeScene  string    `xorm:"not null default '' comment('关注场景') VARCHAR(80)"`
	QrScene         int       `xorm:"not null default 0 comment('二维码扫码场景') INT(11)"`
	QrSceneStr      string    `xorm:"not null default '' comment('二维码扫码场景描述') VARCHAR(128)"`
	Groupid         int       `xorm:"not null default 0 comment('用户所在的分组ID') INT(128)"`
	OpenidPlatformA string    `xorm:"not null default '' comment('boss小程序openid') VARCHAR(128)"`
	OpenidPlatformB string    `xorm:"not null default '' comment('手机app的openid') VARCHAR(128)"`
	CreatedAt       time.Time `xorm:"not null comment('创建时间') DATETIME"`
}
