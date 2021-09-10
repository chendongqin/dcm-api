package es

type DyAwemeSumCount struct {
	AvgDigg struct {
		Value float64 `json:"value"`
	} `json:"avgDigg"`
	TotalSales struct {
		Value float64 `json:"value"`
	} `json:"totalSales"`
	AvgShare struct {
		Value float64 `json:"value"`
	} `json:"avgShare"`
	AvgComment struct {
		Value float64 `json:"value"`
	} `json:"avgComment"`
	TotalGmv struct {
		Value float64 `json:"value"`
	} `json:"totalGmv"`
}
