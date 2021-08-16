package entity

var DyAwemeShareTopMap = HbaseEntity{
	"crawl_time": {Long, "crawl_time"},
	"data":       {AJson, "data"},
}

type DyAwemeShareTops struct {
	CrawlTime int64             `json:"crawl_time"`
	Data      []DyAwemeShareTop `json:"data"`
}

type DyAwemeShareTop struct {
	AuthorId      int64  `json:"author_id"`
	Category      string `json:"category"`
	Avatar        string `json:"avatar"`
	FollowerCount int64  `json:"follower_count"`
	IncDiggCount  int64  `json:"incDiggCount"`
	InitRank      int    `json:"init_rank"`
	Name          string `json:"name"`
	RankChange    int    `json:"rank_change"`
	Score         int64  `json:"score"`
	ShortId       string `json:"short_id"`
	UniqueId      string `json:"unique_id"`
}

type DyAwemeShareTopCopy struct {
	AuthorId      string `json:"author_id"`
	Category      string `json:"category"`
	Avatar        string `json:"avatar"`
	FollowerCount int64  `json:"follower_count"`
	IncDiggCount  int64  `json:"incDiggCount"`
	InitRank      int    `json:"init_rank"`
	Name          string `json:"name"`
	RankChange    int    `json:"rank_change"`
	Score         int64  `json:"score"`
	ShortId       string `json:"short_id"`
	UniqueId      string `json:"unique_id"`
}
