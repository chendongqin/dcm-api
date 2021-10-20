package dcm

type DyAuthorIncome struct {
	Id         int    `xorm:"not null pk autoincr INT(11)"`
	AuthorId   string `xorm:"not null default '' comment('达人id') VARCHAR(50)"`
	UserId     int    `xorm:"not null default 0 comment('用户id') INT(11)"`
	CreateTime int    `xorm:"not null default 0 comment('创建时间') INT(10)"`
}
