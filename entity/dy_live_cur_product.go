package entity

var DyLiveCurProductMap = HbaseEntity{
	"room_id":          {String, "room_id"},
	"room_create_time": {Long, "room_create_time"},
	"crawl_time":       {Long, "crawl_time"},
	"index":            {Int, "index"},
	"promotion":        {AJson, "promotion"},
}

type DyLiveCurProduct struct {
	RoomID         string               `json:"room_id"`
	RoomCreateTime int64                `json:"room_create_time"`
	CrawlTime      int64                `json:"crawl_time"`
	Index          int                  `json:"index"`
	Promotion      []DyLiveCurPromotion `json:"promotion"`
}

type DyLiveCurPromotion struct {
	PromotionID     string  `json:"promotion_id"`      //
	ProductID       string  `json:"product_id"`        //全网id
	StartTime       int64   `json:"start_time"`        //讲解开始时间
	EndTime         int64   `json:"end_time"`          //讲解结束时间
	StartCnt        int64   `json:"start_cnt"`         //讲解开始时间正在去购买人数
	EndCnt          int64   `json:"end_cnt"`           //讲解结束时间正在去购买人数
	StartUserCount  int64   `json:"start_user_count"`  //开始讲解在线人数
	StartTotalUser  int64   `json:"start_total_user"`  //开始讲解在线总pv
	EndUserCount    int64   `json:"end_user_count"`    //
	EndTotalUser    int64   `json:"end_total_user"`    //
	StartSales      int64   `json:"start_sales"`       //
	EndSales        int64   `json:"end_sales"`         //
	StartPrice      float64 `json:"start_price"`       //
	EndPrice        float64 `json:"end_price"`         //
	Pv              int64   `json:"pv"`                //讲解pv
	UserCount       int64   `json:"user_count"`        //讲解在线人数
	Sales           int64   `json:"sales"`             //爬取时全网销量
	TotalUserCount  int64   `json:"total_user_count"`  //
	TotalCrawlTimes int64   `json:"total_crawl_times"` //
	ShopId          string  `json:"shop_id"`           //
	ShopName        string  `json:"shop_name"`         //
	ShopIcon        string  `json:"shop_icon"`         //
	PriceMax        float64 `json:"price_max"`
	PriceMin        float64 `json:"price_min"`
}
