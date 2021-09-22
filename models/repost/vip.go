package repost

type DyTeamSubRet struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	UserVipId     int    `json:"user_vip_id"`
	Remark        string `json:"remark"`
	UpdateTime    int64  `json:"update_time"`
	LoginTime     int64  `json:"login_time"`
	SubExpiration int64  `json:"sub_expiration"`
}
