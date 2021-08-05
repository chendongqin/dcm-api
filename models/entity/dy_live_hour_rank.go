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
	Category string `json:"category"`
	LiveInfo struct {
		Cover      string `json:"cover"`
		CreateTime int    `json:"create_time"`
		Tag        string `json:"tag"`
		Title      string `json:"title"`
		TotalUser  int    `json:"total_user"`
		UserCount  int    `json:"user_count"`
		User       struct {
			Avatar        string `json:"avatar"`
			DisplayId     string `json:"display_id"`
			FollowerCount int    `json:"follower_count"`
			Gender        int    `json:"gender"`
			Id            string `json:"id"`
			Nickname      string `json:"nickname"`
			RoomId        string `json:"room_id"`
			ShortId       string `json:"short_id"`
		} `json:"user"`
	} `json:"live_info"`
	Rank     int    `json:"rank"`
	RoomId   string `json:"room_id"`
	ShareUrl string `json:"share_url"`
	Score    int    `json:"score"`
}
