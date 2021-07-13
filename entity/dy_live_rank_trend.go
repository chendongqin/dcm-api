package entity

var DyLiveRankTrendsMap = HbaseEntity{
	"crawl_time":    {Long, "crawl_time"},
	"discover_time": {Long, "discover_time"},
	"rank_data":     {AJson, "rank_data"},
}

type DyLiveRankTrends struct {
	CrawlTime    int64             `json:"crawl_time"`
	DiscoverTime int64             `json:"discover_time"`
	RankData     []DyLiveRankTrend `json:"rank_data"`
}

type DyLiveRankTrend struct {
	CrawlTime int64 `json:"crawl_time"`
	Rank      int   `json:"rank"`
	Type      int   `json:"type"` //类型 1:小时榜 8: 带货小时榜
}
