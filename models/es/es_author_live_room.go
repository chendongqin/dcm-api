package es

type EsAuthorLiveRoom struct {
	RoomID          string  `json:"room_id"`
	AuthorID        string  `json:"author_id"`
	Title           string  `json:"title"`
	CreateTime      string  `json:"create_time"`
	CreateTimestamp int64   `json:"create_timestamp"`
	MaxUserCount    int     `json:"max_user_count"`
	NumProducts     int     `json:"num_products"`
	PredictSales    float64 `json:"predict_sales"`
	PredictGmv      float64 `json:"predict_gmv"`
	RealSales       float64 `json:"real_sales"`
	RealGmv         float64 `json:"real_gmv"`
	Gmv             float64 `json:"gmv"`
	Sales           float64 `json:"sales"`
	Cover           string  `json:"cover"`
}
