package dcm

type DcWechatMenuClick struct {
	Id   int    `xorm:"not null pk autoincr INT(11)"`
	Name string `xorm:"not null default '' VARCHAR(30)"`
	Key  string `xorm:"not null default '' VARCHAR(30)"`
}
