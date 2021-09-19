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

type EsLiveSumDataCategoryLevel struct {
	Key           string `json:"key"`
	DocCount      int    `json:"doc_count"`
	TotalWatchCnt struct {
		Value int64 `json:"value"`
	} `json:"total_watch_cnt"`
	StatsCustomerUnitPrice struct {
		Count int     `json:"count"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
		Avg   float64 `json:"avg"`
		Sum   float64 `json:"sum"`
	} `json:"stats_customer_unit_price"`
	TotalGmv struct {
		Value float64 `json:"value"`
	} `json:"total_gmv"`
}

type LiveSumDataCategoryLevel struct {
	Level             string  `json:"level"`
	RoomCount         int     `json:"room_count"`
	TotalWatch        int64   `json:"total_watch"`
	AvgWatch          int64   `json:"avg_watch"`
	WatchPer          float64 `json:"watch_per"`
	TotalGmv          float64 `json:"total_gmv"`
	AvgGmv            float64 `json:"avg_gmv"`
	GmvPer            float64 `json:"gmv_per"`
	CustomerUnitPrice struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"customer_unit_price"`
}

type EsLiveSumDataCategoryLevelTwo struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
	LiveTwo  struct {
		Buckets []EsLiveSumDataCategoryLevelTwoBucket `json:"buckets"`
	} `json:"live_tow"`
}

type EsLiveSumDataCategoryLevelTwoBucket struct {
	Key           int `json:"key"`
	DocCount      int `json:"doc_count"`
	TotalWatchCnt struct {
		Value int64 `json:"value"`
	} `json:"total_watch_cnt"`
	TotalGmv struct {
		Value float64 `json:"value"`
	} `json:"total_gmv"`
}

type LiveSumDataCategoryLevelTwo struct {
	FlowLevel     string  `json:"flow_level"`
	StayLevel     int     `json:"stay_level"`
	RoomCount     int     `json:"room_count"`
	TotalWatchCnt int64   `json:"total_watch_cnt"`
	TotalGmv      float64 `json:"total_gmv"`
}
