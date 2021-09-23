package es

type EsAuthorLiveProduct struct {
	RoomID         string  `json:"room_id"`
	AuthorID       string  `json:"author_id"`
	RoomTitle      string  `json:"room_title"`
	RoomCover      string  `json:"room_cover"`
	Title          string  `json:"title"`
	MaxUserCount   int64   `json:"max_user_count"`
	ElasticTitle   string  `json:"elastic_title"`
	Nickname       string  `json:"nickname"`
	Avatar         string  `json:"avatar"`
	Cover          string  `json:"cover"`
	ProductID      string  `json:"product_id"`
	ExtInfo        string  `json:"ext_info"`
	ForSale        int     `json:"for_sale"`
	StartTime      int64   `json:"start_time"`
	ShelfTime      int64   `json:"shelf_time"`
	LiveCreateTime int64   `json:"live_create_time"`
	Pv             int64   `json:"pv"`
	GPM            float64 `json:"gpm"`
	Price          float64 `json:"price"`
	IsReturn       int     `json:"is_return"` //是否返场
	DcmLevelFirst  string  `json:"dcm_level_first"`
	FirstCname     string  `json:"first_cname"`
	SecondCname    string  `json:"second_cname"`
	RoomStatus     string  `json:"room_status"`
	ThirdCname     string  `json:"third_cname"`
	BrandName      string  `json:"brand_name"`
	PlatformLabel  string  `json:"platform_label"`
	ShopName       string  `json:"shop_name"`
	ShopId         string  `json:"shop_id"`
	ShopIcon       string  `json:"shop_icon"`
	CreateTime     string  `json:"create_time"` //直播间时间
	PredictSales   float64 `json:"predict_sales"`
	PredictGmv     float64 `json:"predict_gmv"`
	RealGmv        float64 `json:"real_gmv"`
	BuyRate        float64 `json:"buy_rate"`
	AuthorRoomID   string  `json:"author_room_id"`
}
