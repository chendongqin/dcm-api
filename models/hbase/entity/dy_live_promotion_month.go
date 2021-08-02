package entity

var DyLivePromotionMonthMap = HbaseEntity{
	"dy_promotion_id": {String, "dy_promotion_id"},
	"crawl_time":      {Long, "crawl_time"},
	"monthly":         {AJson, "daily_list"},
	"order_count":     {Long, "order_count"},
	"pv_count":        {Long, "pv_count"},
	"user_count":      {Long, "user_count"},
}

type DyLivePromotionMonth struct {
	DyPromotionID string               `json:"dy_promotion_id"`
	ProductID     string               `json:"product_id"`
	CrawlTime     int64                `json:"crawl_time"`
	DailyList     []DyLiveProductDaily `json:"daily_list"`
	OrderCount    int64                `json:"order_count"`
	PvCount       int64                `json:"pv_count"`
	UserCount     int64                `json:"user_count"`
}

type DyLiveProductDaily struct {
	ProductOrderAccount  int64  `json:"product_order_account"`
	PromotionUserAccount int64  `json:"promotion_user_account"`
	Pv                   int64  `json:"pv"`
	StatisticsTime       string `json:"statistics_time"`
}
