package dcm

import (
	"time"
)

type DcLiveMonitor struct {
	Id               int       `xorm:"not null pk autoincr INT(11)" json:"id"`
	UserId           int       `xorm:"not null default 0 comment('用户ID') index(user_author) index(user_id) INT(11)" json:"user_id"`
	AuthorId         string    `xorm:"not null default '' comment('监控达人ID') index(user_author) VARCHAR(64)" json:"author_id"`
	OpenId           string    `xorm:"not null default '' comment('微信openID') VARCHAR(128)" json:"-"`
	HasNew           int       `xorm:"not null default 0 comment('是否有新的直播') TINYINT(1)" json:"has_new"`
	Source           int       `xorm:"not null default 0 comment('监控来源 0 洞察猫') TINYINT(4)" json:"source"`
	StartTime        time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('预约开始时间') index(status_start_time) TIMESTAMP" json:"-"`
	EndTime          time.Time `xorm:"not null default '0000-00-00 00:00:00' comment('预约结束时间') index(status_end_time) TIMESTAMP" json:"-"`
	NextTime         int64     `xorm:"not null default 0 comment('下一次监控时间') BIGINT(20)" json:"next_time"`
	Status           int       `xorm:"not null default 0 comment('状态 0 等待监控 1监控中 2监控结束 3取消监控') index(status_end_time) index(status_start_time) TINYINT(4)" json:"status"`
	DelStatus        int       `xorm:"not null default 0 comment('是否删除') TINYINT(1)" json:"del_status"`
	Notice           int       `xorm:"not null default 0 comment('开播提醒') TINYINT(1)" json:"notice"`
	FinishNotice     int       `xorm:"not null default 0 comment('下播提醒') TINYINT(1)" json:"finish_notice"`
	ProductId        string    `xorm:"not null default '' comment('商品ID，逗号分割') VARCHAR(512)" json:"product_id"`
	CreateTime       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间') index(user_id) TIMESTAMP" json:"-"`
	UpdateTime       time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' comment('结束时间') TIMESTAMP" json:"-"`
	FreeCount        int       `xorm:"not null default 0 comment('免费次数') INT(11)" json:"free_count"`
	PurchaseCount    int       `xorm:"not null default 0 comment('付费次数') INT(11)" json:"purchase_count"`
	Top              int       `xorm:"not null default 0 comment('是否置顶 0否 1是') TINYINT(1)" json:"top"`
	StartTimeString  string    `xorm:"-" json:"start_time"`
	EndTimeString    string    `xorm:"-" json:"end_time"`
	CreateTimeString string    `xorm:"-" json:"create_time"`
	UpdateTimeString string    `xorm:"-" json:"update_time"`
	RoomId           string    `xorm:"-" json:"room_id"`
	RoomCount        int       `xorm:"-" json:"room_count"`
	Nickname         string    `xorm:"-" json:"nickname"`
	Avatar           string    `xorm:"-" json:"avatar"`
	UniqueID         string    `xorm:"-" json:"unique_id"`
}
