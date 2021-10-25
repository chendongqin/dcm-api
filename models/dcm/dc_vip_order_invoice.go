package dcm

import (
	"time"
)

type DcVipOrderInvoice struct {
	Id          int       `xorm:"not null pk autoincr INT(11)"`
	UserId      int       `xorm:"not null default 0 comment('用户id') INT(11)"`
	Username    string    `xorm:"not null default '' comment('用户名') VARCHAR(32)"`
	Amount      float64   `xorm:"default 0.00 comment('开票金额') DECIMAL(10,2)"`
	Head        string    `xorm:"not null default '' comment('发票抬头') VARCHAR(100)"`
	HeadType    int       `xorm:"not null default 0 comment('抬头类型0企业1个人') TINYINT(1)"`
	InvoiceNum  string    `xorm:"not null default '' comment('企业纳税识别号') VARCHAR(100)"`
	Email       string    `xorm:"not null default '' comment('电子邮箱') VARCHAR(100)"`
	Phone       string    `xorm:"not null default '' comment('手机号') VARCHAR(100)"`
	CompanyTel  string    `xorm:"not null default '' comment('公司电话') VARCHAR(100)"`
	BankName    string    `xorm:"not null default '' comment('开户行') VARCHAR(50)"`
	BankAccount string    `xorm:"not null default '' comment('开户行账号') VARCHAR(50)"`
	RegAddress  string    `xorm:"not null default '' comment('注册地址') VARCHAR(255)"`
	Remark      string    `xorm:"not null default '' comment('发票备注') VARCHAR(100)"`
	InvoiceType int       `xorm:"not null default 0 comment('发票类型') TINYINT(1)"`
	Status      int       `xorm:"not null default 0 comment('0申请中1已通过2已拒绝3已发送') TINYINT(1)"`
	Address     string    `xorm:"not null default '' comment('发票地址') VARCHAR(255)"`
	CreateTime  time.Time `xorm:"TIMESTAMP"`
	UpdateTime  time.Time `xorm:"TIMESTAMP"`
}
