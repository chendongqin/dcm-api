package entity

var DyLiveHourPopularityRanksMap = HbaseEntity{
	"begin_time": {Long, "begin_time"},
	"delta_time": {Long, "delta_time"},
	"crawl_time": {Long, "crawl_time"},
	"ranks":      {AJson, "ranks"},
}

type DyLiveHourPopularityRanks struct {
	Ranks     []DyLiveHourPopularityRank `json:"ranks"`
	CrawlTime int64                      `json:"crawl_time"`
}

type DyLiveHourPopularityRank struct {
	Contributor  int64            `json:"contributor"`
	IncFansNum   int64            `json:"inc_fans_num"`
	LiveInfo     LiveRankLiveInfo `json:"live_info"`
	Rank         int              `json:"rank"`
	RoomId       string           `json:"room_id"`
	Score        int              `json:"score"`
	ShareUrl     string           `json:"share_url"`
	TotalFansNum int64            `json:"total_fans_num"`
}
