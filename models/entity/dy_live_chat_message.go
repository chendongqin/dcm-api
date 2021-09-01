package entity

var DyLiveChatMessageMap = HbaseEntity{
	"latest500_msg": {AJson, "latest500_msg"},
	"end_num":       {Long, "end_num"},
}

type DyLiveChatMessage struct {
	Latest500Msg []LivingChatMessage `json:"latest500_msg"`
	EndNum       int64               `json:"end_num"`
}

type LivingChatMessage struct {
	CreateTime int64  `json:"create_time"`
	Avatar     string `json:"avatar"`
	Content    string `json:"content"`
	Nickname   string `json:"nickname"`
	RankId     int64  `json:"rank_id"`
}
