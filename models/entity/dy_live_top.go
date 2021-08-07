package entity

var DyLiveTopMap = HbaseEntity{
	"begin_time": {Long, "begin_time"},
	"delta_time": {Long, "delta_time"},
	"crawl_time": {Long, "crawl_time"},
	"ranks":      {AJson, "ranks"},
}

type DyLiveTopRanks struct {
	Ranks     []DyLiveRank `json:"ranks"`
	CrawlTime int64        `json:"crawl_time"`
}

type LiveRankLiveInfo struct {
	Cover      string         `json:"cover"`
	CreateTime int            `json:"create_time"`
	Tag        string         `json:"tag"`
	Title      string         `json:"title"`
	TotalUser  int            `json:"total_user"`
	UserCount  int            `json:"user_count"`
	User       LiveRankAuthor `json:"user"`
}

type LiveRankAuthor struct {
	Avatar        string `json:"avatar"`
	DisplayId     string `json:"display_id"`
	FollowerCount int    `json:"follower_count"`
	Gender        int    `json:"gender"`
	Id            string `json:"id"`
	Nickname      string `json:"nickname"`
	ShortId       string `json:"short_id"`
}

type DyLiveRank struct {
	Category string           `json:"category"`
	LiveInfo LiveRankLiveInfo `json:"live_info"`
	Rank     int              `json:"rank"`
	RoomId   string           `json:"room_id"`
	ShareUrl string           `json:"share_url"`
	Score    int              `json:"score"`
}
