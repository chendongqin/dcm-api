package dy

import "time"

type RepostAccountData struct {
	UserId      int                   `json:"user_id"`
	Username    string                `json:"username"`
	Nickname    string                `json:"nickname"`
	Avatar      string                `json:"avatar"`
	PasswordSet int                   `json:"password_set"`
	Wechat      int                   `json:"wechat"`
	BindPhone   int                   `json:"bind_phone"`
	CollectSum  CollectSum            `json:"collect_sum"`
	DyLevel     RepostAccountVipLevel `json:"dy_level"`
	XhsLevel    RepostAccountVipLevel `json:"xhs_level"`
	TbLevel     RepostAccountVipLevel `json:"tb_level"`
}

type CollectSum struct {
	Author  int
	Product int
	Aweme   int
}

type RepostAccountToken struct {
	UserId      int    `json:"user_id"`
	TokenString string `json:"token_string"`
	ExpTime     int64  `json:"exp_time"`
}

type RepostAccountVipLevel struct {
	Level             int    `json:"level"`
	LevelName         string `json:"level_name"`
	ExpirationTime    string `json:"expiration_time"`
	ExpirationDays    int    `json:"expiration_days"`
	SubExpirationDays int    `json:"sub_expiration_days"`
	SubNum            int    `json:"sub_num"`
	IsSub             int    `json:"is_sub"`
	SubExpirationTime string `json:"sub_expiration_time"`
	ParentId          int    `json:"parent_id"`
}

type AccountVipLevel struct {
	Id                int       `json:"id"`
	PlatForm          int       `json:"plat_form"`
	Level             int       `json:"level"`
	FeeLiveMonitor    int       `json:"fee_live_monitor"`
	SubNum            int       `json:"sub_num"`
	IsSub             int       `json:"is_sub"`
	ExpirationTime    time.Time `json:"expiration_time"`
	SubExpirationTime time.Time `json:"sub_expiration_time"`
	ParentId          int       `json:"parent_id"`
}
