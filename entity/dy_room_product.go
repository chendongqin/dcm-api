package entity

var DyRoomProductMap = HbaseEntity{
	"trend_data": {AJson, "trend_data"},
	"price":      {Double, "price"},
	"author_id":  {String, "author_id"},
	"room_id":    {String, "room_id"},
}

type DyRoomProduct struct {
	TrendData []DyRoomProductTrend `json:"trend_data"`
	Price     float64              `json:"price"`
	AuthorId  string               `json:"author_id"`
	RoomId    string               `json:"room_id"`
}

type DyRoomProductTrend struct {
	CrawlTime int64   `json:"crawl_time"`
	Price     float64 `json:"price"`
	Sales     float64 `json:"sales"`
}
