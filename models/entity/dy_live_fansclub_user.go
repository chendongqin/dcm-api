package entity

var DyLiveFansClubUserMap = HbaseEntity{
	"fans_infos": {AJson, "fans_infos"},
}

type DyLiveFansClubUser struct {
	FansInfos []DyLiveFansClubUserInfo `json:"fans_infos"`
}

type DyLiveFansClubUserInfo struct {
	Intimacy int64  `json:"intimacy"`
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
