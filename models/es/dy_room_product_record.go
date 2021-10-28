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
		Value int64 `json:"value"`
	} `json:"live_create_time"`
	PredictGmv struct {
		Value float64 `json:"value"`
	} `json:"predict_gmv"`
	PredictSales struct {
		Value float64 `json:"value"`
	} `json:"predict_sales"`
}

type LiveAuthorProduct struct {
	RoomId            string  `json:"room_id"`
	RoomStatus        string  `json:"room_status"`
	AuthorId          string  `json:"author_id"`
	Title             string  `json:"title"`
	Cover             string  `json:"cover"`
	ProductId         string  `json:"product_id"`
	ExtInfo           string  `json:"ext_info"`
	ForSale           int     `json:"for_sale"`
	ShelfTime         int64   `json:"shelf_time"`
	StartTime         int64   `json:"start_time"`
	Pv                int64   `json:"pv"`
	Price             float64 `json:"price"`
	IsReturn          int     `json:"is_return"`
	DcmLevelFirst     string  `json:"dcm_level_first"`
	FirstCname        string  `json:"first_cname"`
	SecondCname       string  `json:"second_cname"`
	ThirdCname        string  `json:"third_cname"`
	PredictSales      float64 `json:"predict_sales"`
	PredictGmv        float64 `json:"predict_gmv"`
	RealSales         float64 `json:"real_sales"`
	RealGmv           float64 `json:"real_gmv"`
	RoomTitle         string  `json:"room_title"`
	Nickname          string  `json:"nickname"`
	MaxUserCount      int64   `json:"max_user_count"`
	CrawlTime         int64   `json:"crawl_time"`
	LiveCreateTime    int64   `json:"live_create_time"`
	Avatar            string  `json:"avatar"`
	RoomCover         string  `json:"room_cover"`
	ElasticTitle      string  `json:"elastic_title"`
	Gpm               float64 `json:"gpm"`
	Dt                string  `json:"dt"`
	ShopId            string  `json:"shop_id"`
	ShopName          string  `json:"shop_name"`
	ShopIcon          string  `json:"shop_icon"`
	BrandName         string  `json:"brand_name"`
	PlatformLabel     string  `json:"platform_label"`
	Tags              string  `json:"tags"`
	Score             float64 `json:"score"`
	Level             int     `json:"level"`
	FinishTime        int64   `json:"finish_time"`
	TotalUser         int64   `json:"total_user"`
	DisplayId         string  `json:"displayId"`
	ShortId           string  `json:"shortId"`
	FollowerCount     int64   `json:"follower_count"`
	AuthorFirstCname  string  `json:"author_first_cname"`
	AuthorSecondCname string  `json:"author_second_cname"`
	FlowRatesIndex    string  `json:"flow_rates_index"`
	FlowRates         string  `json:"flow_rates"`
}

type SumLiveProductAuthor struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
	Data     struct {
		Hits struct {
			Total int `json:"total"`
			Hits  []struct {
				Source LiveAuthorProduct `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	} `json:"data"`
	PredictSales struct {
		Value float64 `json:"value"`
	} `json:"predict_sales"`
	PredictGmv struct {
		Value float64 `json:"value"`
	} `json:"predict_gmv"`
}

type CountAuthorProductRoom struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
	Rooms    struct {
		Buckets []interface{} `json:"buckets"`
	} `json:"rooms"`
	Products struct {
		Buckets []interface{} `json:"buckets"`
	} `json:"products"`
}
