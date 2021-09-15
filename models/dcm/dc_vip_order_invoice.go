package dcm

import (
	"time"
)

type DcVipOrderInvoice struct {
	Id         int       `xorm:"not null pk autoincr INT(11)"`
	UserId     int       `xorm:"not null default 0 comment('用户id') INT(11)"`
	Username   string    `xorm:"not null default '' comment('用户名') VARCHAR(32)"`
	Amount     string    `xorm:"default 0.00 comment('开票金额') DECIMAL(10,2)"`
	Head       string    `xorm:"not null default '' comment('发票抬头') VARCHAR(100)"`
	InvoiceNum string    `xorm:"not null default '' comment('企业纳税识别号') VARCHAR(100)"`
	Emial      string    `xorm:"not null default '' comment('电子邮箱') VARCHAR(100)"`
	Phone      string    `xorm:"not null default '' comment('手机号') VARCHAR(100)"`
	Remark     string    `xorm:"not null default '' comment('发票备注') VARCHAR(100)"`
	Status     int       `xorm:"not null default 0 comment('1申请中2已通过3已拒绝4已发送') TINYINT(1)"`
	Address    string    `xorm:"not null default '' comment('发票地址') VARCHAR(255)"`
	CreateTime time.Time `xorm:"TIMESTAMP"`
	UpdateTime time.Time `xorm:"TIMESTAMP"`
}
