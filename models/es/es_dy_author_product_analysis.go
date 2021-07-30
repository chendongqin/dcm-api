package es

type EsDyAuthorProductAnalysis struct {
	AuthorId          string  `json:"author_id"`
	ProductId         string  `json:"product_id"`
	Price             float64 `json:"price"`
	ShopId            string  `json:"shop_id"`
	ShopName          string  `json:"shop_name"`
	ShopIcon          string  `json:"shop_icon"`
	BrandName         string  `json:"brand_name"`
	Platform          string  `json:"platform"`
	DcmLevelFirst     string  `json:"dcm_level_first"`
	FirstCname        string  `json:"first_cname"`
	SecondCname       string  `json:"second_cname"`
	ThirdCname        string  `json:"third_cname"`
	LivePredictSales  float64 `json:"live_predict_sales"`
	LivePredictGmv    float64 `json:"live_predict_gmv"`
	RoomProductSales  string  `json:"room_product_sales"`
	RoomCount         int     `json:"room_count"`
	VedioPredictSales float64 `json:"vedio_predict_sales"`
	VideoPredictGmv   float64 `json:"video_predict_gmv"`
	VedioProductSales float64 `json:"vedio_product_sales"`
	VedioCount        int     `json:"vedio_count"`
	RowTime           string  `json:"row_time"`
	CreateTime        string  `json:"create_time"`
	ShelfTime         int64   `json:"shelf_time"`
}
