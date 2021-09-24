package es

type EsGroupByData struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}

type DyAwemeSumCount struct {
	AvgDigg struct {
		Value float64 `json:"value"`
	} `json:"avg_digg"`
	TotalSales struct {
		Value float64 `json:"value"`
	} `json:"total_sales"`
	AvgShare struct {
		Value float64 `json:"value"`
	} `json:"avg_share"`
	AvgComment struct {
		Value float64 `json:"value"`
	} `json:"avg_comment"`
	TotalGmv struct {
		Value float64 `json:"value"`
	} `json:"total_gmv"`
}

type DyLiveSumCount struct {
	Key        string     `json:"key"`
	TotalSales EsSumStats `json:"total_sales"`
	TotalGmv   EsSumStats `json:"total_gmv"`
}

type EsSumStats struct {
	Count int     `json:"count"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Avg   float64 `json:"avg"`
	Sum   float64 `json:"sum"`
}

type DyLiveDataUserSumCount struct {
	TotalWatchCnt struct {
		Value float64 `json:"value"`
	} `json:"total_watch_cnt"`
	TotalUserCount struct {
		Value float64 `json:"value"`
	} `json:"total_user_count"`
}

type DyLiveDataCategorySumCount struct {
	TotalWatchCnt struct {
		Value float64 `json:"value"`
	} `json:"total_watch_cnt"`
	TotalUserCount struct {
		Value float64 `json:"value"`
	} `json:"total_user_count"`
	TotalGmv struct {
		Value float64 `json:"value"`
	} `json:"total_gmv"`
	TotalSales struct {
		Value float64 `json:"value"`
	} `json:"total_sales"`
	TotalTicketCount struct {
		Value float64 `json:"value"`
	} `json:"total_ticket_count"`
}

type DyLiveCategoryRateByWatchCnt struct {
	DocCount int `json:"doc_count"`
	Key      struct {
		DcmLevelFirst string `json:"dcm_level_first"`
	} `json:"key"`
	TotalWatchCnt struct {
		Value int64 `json:"value"`
	} `json:"total_watch_cnt"`
}

type DyLiveCategoryRateByGmv struct {
	DocCount int `json:"doc_count"`
	Key      struct {
		DcmLevelFirst string `json:"dcm_level_first"`
	} `json:"key"`
	TotalGmv struct {
		Value float64 `json:"value"`
	} `json:"total_gmv"`
}

type DyProductAwemeSum struct {
	DocCount int    `json:"doc_count"`
	Key      string `json:"key"`
	TotalGmv struct {
		Value float64 `json:"value"`
	} `json:"total_gmv"`
	TotalSales struct {
		Value int64 `json:"value"`
	} `json:"total_sales"`
}
