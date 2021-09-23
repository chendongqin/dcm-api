package dcm

import (
	"time"
)

type DcUserDyCollect struct {
	Id          int       `xorm:"not null pk autoincr INT(11)"`
	UserId      int       `xorm:"not null default 0 comment('用户id') index(UN_INDEX) INT(11)"`
	CollectType int       `xorm:"not null default 1 comment('1达人2商品3视频4小店') index(UN_INDEX) TINYINT(1)"`
	CollectId   string    `xorm:"not null default '' comment('type:1达人id2商品id3视频id4小店id') index(UN_INDEX) VARCHAR(50)"`
	Label       string    `xorm:"not null default '' comment('类型') VARCHAR(20)"`
	UniqueId    string    `xorm:"not null default '' comment('抖音号') VARCHAR(20)"`
	Nickname    string    `xorm:"not null default '' comment('达人昵称') VARCHAR(50)"`
	PromotionID string    `xorm:"not null default '' comment('抖音商品id') VARCHAR(50)"`
	TagId       int       `xorm:"not null default 0 comment('分组id') INT(11)"`
	Remark      string    `xorm:"not null default '' comment('备注') VARCHAR(100)"`
	Status      int       `xorm:"not null default 1 comment('状态0取消1正常') TINYINT(1)"`
	CreateTime  time.Time `xorm:"TIMESTAMP"`
	UpdateTime  time.Time `xorm:"TIMESTAMP"`
}
