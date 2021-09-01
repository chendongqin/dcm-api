package entity

var DyLiveChatMessageMap = HbaseEntity{
	"latest500_msg": {AJson, "latest500_msg"},
	"end_num":       {Long, "end_num"},
	"visit_num":     {Long, "visit_num"},
	"visits":        {AJson, "visits"},
}

type DyLiveChatMessage struct {
	Latest500Msg []LivingChatMessage `json:"latest500_msg"`
	Visits       []LivingChatVisit   `json:"visits"`
	EndNum       int64               `json:"end_num"`
	VisitNum     int64               `json:"visit_num"`
}

type LivingChatMessage struct {
	CreateTime int64  `json:"create_time"`
	Avatar     string `json:"avatar"`
	Content    string `json:"content"`
	Nickname   string `json:"nickname"`
	RankId     int64  `json:"rank_id"`
}

type LivingChatVisit struct {
	RankId   int64  `json:"rank_id"`
	Nickname string `json:"nickname"`
}
