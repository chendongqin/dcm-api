package dcm

type DcWechatMenuClick struct {
	Id          int    `xorm:"not null pk autoincr INT(11)"`
	Name        string `xorm:"not null default '' VARCHAR(30)"`
	MsgKey      string `xorm:"not null default '' VARCHAR(30)"`
	Type        string `xorm:"not null default '' VARCHAR(10)"`
	MediaId     string `xorm:"not null default '' VARCHAR(100)"`
	Content     string `xorm:"not null default '' VARCHAR(255)"`
	Title       string `xorm:"not null default '' VARCHAR(100)"`
	Description string `xorm:"not null default '' VARCHAR(255)"`
}
