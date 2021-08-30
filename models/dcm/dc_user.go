package dcm

import (
	"time"
)

type DcUser struct {
	Id               int       `xorm:"not null pk autoincr INT(11)"`
	Username         string    `xorm:"not null default '' comment('用户名（手机号）') CHAR(11)"`
	Nickname         string    `xorm:"not null comment('昵称') VARCHAR(30)"`
	Password         string    `xorm:"not null default '' comment('密码') VARCHAR(32)"`
	Salt             string    `xorm:"not null comment('盐') CHAR(4)"`
	Avatar           string    `xorm:"not null default '' comment('头像') VARCHAR(255)"`
	Successions      int       `xorm:"not null default 0 comment('连续登陆天数') SMALLINT(5)"`
	MaxSuccessions   int       `xorm:"not null default 0 comment('最大连续登陆天数') SMALLINT(5)"`
	TotalSuccessions int       `xorm:"not null default 0 comment('总的登陆天数') SMALLINT(5)"`
	PrevTime         time.Time `xorm:"comment('最近活跃时间') TIMESTAMP"`
	LoginTime        time.Time `xorm:"comment('登陆时间') TIMESTAMP"`
	LoginIp          string    `xorm:"not null default '' comment('登陆ip') VARCHAR(50)"`
	Status           int       `xorm:"not null default 1 comment('状态1正常0禁用') TINYINT(1)"`
	Openid           string    `xorm:"not null default '' comment('openid') VARCHAR(100)"`
	Unionid          string    `xorm:"not null default '' comment('unionid') VARCHAR(100)"`
	OpenidApp        string    `xorm:"not null default '' comment('客户端openid') VARCHAR(100)"`
	CreateTime       time.Time `xorm:"comment('创建时间') TIMESTAMP"`
	UpdateTime       time.Time `xorm:"comment('更新时间') TIMESTAMP"`
	SetPassword      int       `xorm:"not null default 0 comment('是否设置了登陆密码') TINYINT(1)"`
	Entrance         int       `xorm:"not null comment('用户来源0:PC,1:小程序,2:APP,3:wap') TINYINT(1)"`
	IsInstallApp     int       `xorm:"not null comment('是否安装app') TINYINT(1)"`
}
