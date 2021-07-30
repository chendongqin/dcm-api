package dcm

type DcProductCate struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"`
	Name     string `xorm:"not null default '' comment('名称') VARCHAR(50)"`
	Level    int    `xorm:"not null comment('级别') TINYINT(1)"`
	ParentId int    `xorm:"not null default 0 comment('父id') INT(11)"`
}
