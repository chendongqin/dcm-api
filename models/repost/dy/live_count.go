package dy

type LiveSumCountByCategoryBase struct {
	RoomNum   int     `json:"room_num"`
	WatchCnt  int64   `json:"watch_cnt"`
	UserCount int64   `json:"user_count"`
	Gmv       float64 `json:"gmv"`
	BuyRate   float64 `json:"buy_rate"`
	Uv        float64 `json:"uv"`
}

type LiveSumCountByCategory struct {
	RoomNum           int     `json:"room_num"`
	WatchCnt          int64   `json:"watch_cnt"`
	UserCount         int64   `json:"user_count"`
	Gmv               float64 `json:"gmv"`
	BuyRate           float64 `json:"buy_rate"`
	Uv                float64 `json:"uv"`
	RoomNumMonthInc   float64 `json:"room_num_month_inc"`
	RoomNumLastInc    float64 `json:"room_num_last_inc"`
	WatchCntMonthInc  float64 `json:"watch_cnt_month_inc"`
	WatchCntLastInc   float64 `json:"watch_cnt_last_inc"`
	UserCountMonthInc float64 `json:"user_count_month_inc"`
	UserCountLastInc  float64 `json:"user_count_last_inc"`
	GmvMonthInc       float64 `json:"gmv_month_inc"`
	GmvLastInc        float64 `json:"gmv_last_inc"`
	BuyRateMonthInc   float64 `json:"buy_rate_month_inc"`
	BuyRateLastInc    float64 `json:"buy_rate_last_inc"`
	UvMonthInc        float64 `json:"uv_month_inc"`
	UvLastInc         float64 `json:"uv_last_inc"`
}
