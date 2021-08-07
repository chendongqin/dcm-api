package entity

var DyLiveHourSellRanksMap = HbaseEntity{
	"begin_time": {Long, "begin_time"},
	"delta_time": {Long, "delta_time"},
	"crawl_time": {Long, "crawl_time"},
	"ranks":      {AJson, "ranks"},
}

type DyLiveHourSellRanks struct {
	Ranks     []DyLiveHourSellRank `json:"ranks"`
	CrawlTime int64                `json:"crawl_time"`
}

type DyLiveHourSellRank struct {
	LiveInfo     LiveRankLiveInfo `json:"live_info"`
	PredictGmv   float64          `json:"predict_gmv"`
	PredictSales float64          `json:"predict_sales"`
	Rank         int              `json:"rank"`
	RealGmv      float64          `json:"real_gmv"`
	RealSales    float64          `json:"real_sales"`
	RoomId       string           `json:"room_id"`
	Score        int              `json:"score"`
	ShareUrl     string           `json:"share_url"`
	ShopTags     []string         `json:"shop_tags"`
}
