package entity

var DyLiveFansClubMap = HbaseEntity{
	"crawl_time":           {Long, "crawl_time"},
	"author_id":            {String, "author_id"},
	"room_id":              {String, "room_id"},
	"active_fans_count":    {Long, "active_fans_count"},
	"hot_rank":             {Int, "hot_rank"},
	"name":                 {String, "name"},
	"today_new_fans_count": {Long, "today_new_fans_count"},
	"total_fans_count":     {Long, "total_fans_count"},
	"discount_price":       {Int, "discount_price"},
}

type DyLiveFansClub struct {
	CrawlTime         int64  `json:"crawl_time"`
	AuthorID          string `json:"author_id"`
	RoomID            string `json:"room_id"`
	ActiveFansCount   int64  `json:"active_fans_count"`
	HotRank           int    `json:"hot_rank"`
	Name              string `json:"name"`
	TodayNewFansCount int64  `json:"today_new_fans_count"`
	TotalFansCount    int64  `json:"total_fans_count"`
	DiscountPrice     int    `json:"discount_price"`
}
