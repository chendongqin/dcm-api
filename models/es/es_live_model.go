package es

type EsAuthorLiveRoom struct {
	RoomID          string  `json:"room_id"`
	AuthorID        string  `json:"author_id"`
	Title           string  `json:"title"`
	CreateTime      string  `json:"create_time"`
	CreateTimestamp int64   `json:"create_timestamp"`
	MaxUserCount    int     `json:"max_user_count"`
	NumProducts     int     `json:"num_products"`
	Sales           float64 `json:"sales"`
	Gmv             float64 `json:"gmv"`
	Cover           string  `json:"cover"`
}

type EsAuthorLiveProduct struct {
	RoomID        string  `json:"room_id"`
	AuthorID      string  `json:"author_id"`
	Title         string  `json:"title"`
	Cover         string  `json:"cover"`
	ProductID     string  `json:"product_id"`
	ExtInfo       string  `json:"ext_info"`
	ForSale       int     `json:"for_sale"`
	StartTime     int64   `json:"start_time"`
	ShelfTime     int64   `json:"shelf_time"`
	Pv            int64   `json:"pv"`
	Price         float64 `json:"price"`
	IsReturn      int     `json:"is_return"` //是否返场
	DcmLevelFirst string  `json:"dcm_level_first"`
	FirstCname    string  `json:"first_cname"`
	SecondCname   string  `json:"second_cname"`
	//ThirdCname    string  `json:"third_cname"`
	CreateTime   string  `json:"create_time"` //直播间时间
	PredictSales float64 `json:"predict_sales"`
	PredictGmv   float64 `json:"predict_gmv"`
	RealGmv      float64 `json:"real_gmv"`
	BuyRate      float64 `json:"buy_rate"`
}
