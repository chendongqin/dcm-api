package dy

type LivingInfo struct {
	RoomId         string  `json:"room_id"`
	AuthorId       string  `json:"author_id"`
	Title          string  `json:"title"`
	Cover          string  `json:"cover"`
	CreateTime     int64   `json:"create_time"`
	Gmv            float64 `json:"gmv"`
	UserCount      int64   `json:"user_count"`
	TotalUserCount int64   `json:"total_user_count"`
	RoomStatus     int     `json:"room_status"`
	FinishTime     int64   `json:"finish_time"`
	LiveTime       int64   `json:"live_time"`
	Uv             float64 `json:"uv"`
	AvgOnlineTime  float64 `json:"avg_online_time"`
	BarrageRate    float64 `json:"barrage_rate"`
}

type LivingProducts struct {
	RoomID        string           `json:"room_id"`
	AuthorID      string           `json:"author_id"`
	RoomTitle     string           `json:"room_title"`
	RoomCover     string           `json:"room_cover"`
	Title         string           `json:"title"`
	ElasticTitle  string           `json:"elastic_title"`
	Cover         string           `json:"cover"`
	ProductID     string           `json:"product_id"`
	ExtInfo       string           `json:"ext_info"`
	ForSale       int              `json:"for_sale"`
	ShelfTime     int64            `json:"shelf_time"`
	Pv            int64            `json:"pv"`
	Price         float64          `json:"price"`
	CreateTime    string           `json:"create_time"`
	StartTime     int64            `json:"start_time"`
	IsReturn      int              `json:"is_return"` //是否返场
	PredictSales  float64          `json:"predict_sales"`
	PredictGmv    float64          `json:"predict_gmv"`
	BuyRate       float64          `json:"buy_rate"`
	CurList       []LiveCurProduct `json:"cur_list"`
	StartCurTime  int64            `json:"start_cur_time"`
	EndCurTime    int64            `json:"end_cur_time"`
	StartPmtSales int64            `json:"start_pmt_sales"`
	EndPmtSales   int64            `json:"end_pmt_sales"`
	CurSecond     int64            `json:"cur_second"`
}
