package dcm

import (
	"time"
)

type DcDirtyLog struct {
	Id           int       `xorm:"not null pk autoincr INT(11)"`
	AdminId      int       `xorm:"not null comment('后台管理员id') INT(11)"`
	ChangeType   int       `xorm:"not null default 1 comment('1达人修改2商品修改3gmv修改') TINYINT(1)"`
	DataId       string    `xorm:"not null default '' comment('操作id：达人id|商品id|直播间id') index VARCHAR(60)"`
	OriginalData string    `xorm:"not null comment('原始数据json') VARCHAR(500)"`
	TargetData   string    `xorm:"not null comment('目标数据') VARCHAR(500)"`
	Status       int       `xorm:"not null default 1 comment('状态1未处理2成功2失败') TINYINT(1)"`
	CreateTime   time.Time `xorm:"TIMESTAMP"`
	UpdateTime   time.Time `xorm:"TIMESTAMP"`
	Remark       string    `xorm:"not null default '' VARCHAR(100)"`
}
