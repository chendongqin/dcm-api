package repost

import (
	"time"
)

type DyTeamSubRet struct {
	Id            int       `json:"id"`
	Username      string    `json:"username"`
	UserVipId     int       `json:"user_vip_id"`
	Remark        string    `json:"remark"`
	UpdateTime    time.Time `json:"update_time"`
	LoginTime     time.Time `json:"login_time"`
	SubExpiration time.Time `json:"sub_expiration"`
}
