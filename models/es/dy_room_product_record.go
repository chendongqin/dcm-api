package es

type LiveAuthorProductList struct {
	Data struct {
		Hits struct {
			Hits []struct {
				Id     string            `json:"_id"`
				Index  string            `json:"_index"`
				Score  interface{}       `json:"_score"`
				Source LiveAuthorProduct `json:"_source"`
				Type   string            `json:"_type"`
				Sort   []int             `json:"sort"`
			} `json:"hits"`
			MaxScore interface{} `json:"max_score"`
			Total    int         `json:"total"`
		} `json:"hits"`
	} `json:"data"`
	DocCount       int    `json:"doc_count"`
	Key            string `json:"key"`
	LiveCreateTime struct {
		Value int `json:"value"`
	} `json:"live_create_time"`
	PredictGmv struct {
		Value float64 `json:"value"`
	} `json:"predict_gmv"`
	PredictSales struct {
		Value int `json:"value"`
	} `json:"predict_sales"`
}

type LiveAuthorProduct struct {
	AuthorId       string  `json:"author_id"`
	Avatar         string  `json:"avatar"`
	BrandName      string  `json:"brand_name"`
	Cover          string  `json:"cover"`
	CrawlTime      int     `json:"crawl_time"`
	DcmLevelFirst  string  `json:"dcm_level_first"`
	DisplayId      string  `json:"displayId,omitempty"`
	Dt             string  `json:"dt"`
	ElasticTitle   string  `json:"elastic_title"`
	ExtInfo        string  `json:"ext_info"`
	FinishTime     int     `json:"finish_time,omitempty"`
	FirstCname     string  `json:"first_cname"`
	FlowRates      string  `json:"flow_rates,omitempty"`
	FlowRatesIndex string  `json:"flow_rates_index,omitempty"`
	FollowerCount  int     `json:"follower_count,omitempty"`
	ForSale        int     `json:"for_sale"`
	Gpm            float64 `json:"gpm"`
	IsReturn       int     `json:"is_return"`
	Level          int     `json:"level,omitempty"`
	LiveCreateTime int     `json:"live_create_time"`
	MaxUserCount   int     `json:"max_user_count"`
	Nickname       string  `json:"nickname"`
	PlatformLabel  string  `json:"platform_label"`
	PredictGmv     float64 `json:"predict_gmv"`
	PredictSales   int     `json:"predict_sales"`
	Price          float64 `json:"price"`
	ProductId      string  `json:"product_id"`
	Pv             int     `json:"pv"`
	RealGmv        float64 `json:"real_gmv"`
	RealSales      int     `json:"real_sales"`
	RoomCover      string  `json:"room_cover"`
	RoomId         string  `json:"room_id"`
	RoomProductId  string  `json:"room_product_id"`
	RoomStatus     string  `json:"room_status"`
	RoomTitle      string  `json:"room_title"`
	Score          float64 `json:"score,omitempty"`
	SecondCname    string  `json:"second_cname"`
	ShelfTime      int     `json:"shelf_time"`
	ShopIcon       string  `json:"shop_icon"`
	ShopId         string  `json:"shop_id"`
	ShopName       string  `json:"shop_name"`
	ShortId        string  `json:"shortId,omitempty"`
	StartTime      int     `json:"start_time"`
	Tags           string  `json:"tags,omitempty"`
	ThirdCname     string  `json:"third_cname"`
	Title          string  `json:"title"`
	TotalUser      int     `json:"total_user,omitempty"`
}
