package entity

var DyLivePmtMap = HbaseEntity{
	"room_status":  {Int, "room_status"},
	"author_id":    {String, "author_id"},
	"room_id":      {String, "room_id"},
	"create_time":  {Long, "create_time"},
	"crawl_time":   {Long, "crawl_time"},
	"purchase_cnt": {Long, "purchase_cnt"},
	"cur":          {String, "cur"},
	"promotions":   {AJson, "promotions"},
	"top":          {Int, "top"},
	"is_bubble":    {Bool, "is_bubble"},
}

type DyLivePmt struct {
	RoomStatus  int               `json:"room_status"` //
	AuthorID    string            `json:"author_id"`
	RoomID      string            `json:"room_id"`
	CreateTime  int               `json:"create_time"`
	CrawlTime   int               `json:"crawl_time"`
	PurchaseCnt int               `json:"purchase_cnt"`
	Cur         string            `json:"cur"`
	Promotions  []DyLivePromotion `json:"promotions"`
	Top         int               `json:"top"`
	IsBubble    bool              `json:"is_bubble"`
}

type DyLivePromotion struct {
	DyPromotionID  string  `json:"dy_promotion_id"`
	ProductID      string  `json:"product_id"`
	ForSale        int     `json:"for_sale"`
	StartTime      int64   `json:"start_time"`
	StopTime       int64   `json:"stop_time"`
	StartUserCount int64   `json:"start_user_count"`
	StartTotalUser int64   `json:"start_total_user"`
	EndUserCount   int64   `json:"end_user_count"`
	EndTotalUser   int64   `json:"end_total_user"`
	Price          float64 `json:"price"`
	//Prices         []float64                   `json:"prices"`
	//PriceTrend     []DyLivePromotionPriceTrend `json:"price_trend"`
	Coupon         float64 `json:"coupon"`
	CouponHeader   string  `json:"coupon_header"`
	InitialSales   int64   `json:"initial_sales"`
	FinalSales     int64   `json:"final_sales"`
	Sales          int64   `json:"sales"`
	InStock        bool    `json:"in_stock"`
	Title          string  `json:"title"`
	ElasticTitle   string  `json:"elastic_title"`
	CosRatio       float64 `json:"cos_ratio"`
	Source         string  `json:"source"`
	ExtInfo        string  `json:"ext_info"`
	Cover          string  `json:"cover"`
	InitialPv      int64   `json:"initial_pv"`
	FinalPv        int64   `json:"final_pv"`
	Pv             int64   `json:"pv"`
	Campaign       bool    `json:"campaign"`
	Index          int     `json:"index"`
	ShopID         string  `json:"shop_id"`
	FlushBuy       bool    `json:"flush_buy"`
	BubbleDuration int     `json:"bubble_duration"`
	BubblePv       int     `json:"bubble_pv"`
	HasH5PmtInfo   bool    `json:"has_h5_pmt_info"`
}

type DyLivePromotionPriceTrend struct {
	CrawlTime int `json:"crawl_time"`
	Price     int `json:"price"`
}
