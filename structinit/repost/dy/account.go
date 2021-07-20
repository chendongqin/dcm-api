package dy

type RepostAccountData struct {
	UserId      int                   `json:"user_id"`
	Username    string                `json:"username"`
	Nickname    string                `json:"nickname"`
	Avatar      string                `json:"avatar"`
	PasswordSet int                   `json:"password_set"`
	DyLevel     RepostAccountVipLevel `json:"dy_level"`
	XhsLevel    RepostAccountVipLevel `json:"xhs_level"`
	TbLevel     RepostAccountVipLevel `json:"tb_level"`
}

type RepostAccountToken struct {
	UserId      int    `json:"user_id"`
	TokenString string `json:"token_string"`
	ExpTime     int64  `json:"exp_time"`
}

type RepostAccountVipLevel struct {
	Level     int    `json:"level"`
	LevelName string `json:"level_name"`
}
