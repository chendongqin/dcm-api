package dcm

type DcAppid struct {
	Id     int    `xorm:"not null pk INT(11)"`
	AppId  int    `xorm:"not null SMALLINT(5)"`
	Secret string `xorm:"not null VARCHAR(16)"`
	Remark string `xorm:"not null default '' VARCHAR(50)"`
}
