package entity

var DyLiveShareTopMap = HbaseEntity{
	"crawl_time": {Long, "crawl_time"},
	"data":       {AJson, "data"},
}

type DyLiveShareTops struct {
	CrawlTime int64            `json:"crawl_time"`
	Data      []DyLiveShareTop `json:"data"`
}

type DyLiveShareTop struct {
	AuthorId   int64                      `json:"author_id"`
	Avatar     string                     `json:"avatar"`
	Category   string                     `json:"category"`
	InitRank   int                        `json:"init_rank"`
	Name       string                     `json:"name"`
	RankChange int                        `json:"rank_change"`
	Score      int64                      `json:"score"`
	ShortId    string                     `json:"short_id"`
	UniqueId   string                     `json:"unique_id"`
	Rooms      map[string]DyLiveShareRoom `json:"rooms"`
}

type DyLiveShareWeekData struct {
	AuthorId   string  `json:"author_id"`
	Avatar     string  `json:"avatar"`
	Category   string  `json:"category"`
	InitRank   int     `json:"init_rank"`
	Name       string  `json:"name"`
	RankChange int     `json:"rank_change"`
	Score      int64   `json:"score"`
	UniqueId   string  `json:"unique_id"`
	Gmv        float64 `json:"gmv"`
	Sales      int64   `json:"sales"`
	TotalUser  int64   `json:"total_user"`
	RoomNum    int     `json:"room_num"`
}

type DyLiveShareRoom struct {
	RoomId       string  `json:"room_id"`
	CreateTime   int64   `json:"create_time"`
	TotalUser    int64   `json:"total_user"`
	PredictGmv   float64 `json:"predict_gmv"`
	PredictSales int64   `json:"predict_sales"`
	RealGmv      float64 `json:"real_gmv"`
	RealSales    int64   `json:"real_sales"`
}
