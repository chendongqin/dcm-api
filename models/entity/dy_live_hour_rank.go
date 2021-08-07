package entity

var DyLiveHourRankMap = HbaseEntity{
	"begin_time": {Long, "begin_time"},
	"delta_time": {Long, "delta_time"},
	"crawl_time": {Long, "crawl_time"},
	"ranks":      {AJson, "ranks"},
}

type DyLiveHourRanks struct {
	Ranks     []DyLiveHourRank `json:"ranks"`
	CrawlTime int64            `json:"crawl_time"`
}

type DyLiveHourRank struct {
	Category string           `json:"category"`
	LiveInfo LiveRankLiveInfo `json:"live_info"`
	Rank     int              `json:"rank"`
	RoomId   string           `json:"room_id"`
	ShareUrl string           `json:"share_url"`
	Score    int              `json:"score"`
}
