package entity

var DyLiveReputationMap = HbaseEntity{
	"type":              {String, "type"},
	"author_reputation": {AJson, "author_reputation"},
	"create_time":       {Long, "create_time"},
	"crawl_time":        {Long, "crawl_time"},
}

type DyLiveReputation struct {
	Type             string                 `json:"type"`
	AuthorReputation DyLiveAuthorReputation `json:"author_reputation"`
	RoomId           string                 `json:"room_id"`
	CreateTime       int64                  `json:"create_time"`
	CrawlTime        int64                  `json:"crawl_time"`
}

type DyLiveAuthorReputation struct {
	Uid   string  `json:"uid"`
	Level int     `json:"level"`
	Score float64 `json:"score"`
}
