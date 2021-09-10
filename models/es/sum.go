package es

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
