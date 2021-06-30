package dy

type AccountData struct {
	UserId   int              `json:"user_id"`
	Username string           `json:"username"`
	Nickname string           `json:"nickname"`
	Avatar   string           `json:"avatar"`
	DyLevel  AccountVipLevel `json:"dy_level"`
	XhsLevel AccountVipLevel `json:"xhs_level"`
	TbLevel  AccountVipLevel `json:"tb_level"`
}

type AccountToken struct {
	UserId      int    `json:"user_id"`
	TokenString string `json:"token_string"`
	ExpTime     int64  `json:"exp_time"`
}

type AccountVipLevel struct {
	Level     int    `json:"level"`
	LevelName string `json:"level_name"`
}
