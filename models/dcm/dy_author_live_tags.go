package dcm

type DyAuthorLiveTags struct {
	Id     int    `xorm:"not null pk autoincr INT(11)"`
	Name   string `xorm:"not null VARCHAR(30)"`
	Weight int    `xorm:"SMALLINT(6)"`
}
