package dy

//销量/销售趋势图
type ShopSaleChart struct {
	Date       []string  `json:"date"`
	SalesCount []int64   `json:"sales_count"`
	GmvCount   []float64 `json:"gmv_count"`
}

//视频/直播趋势图
type ShopLiveAwemeChart struct {
	Date       []string `json:"date"`
	LiveCount  []int64  `json:"live_count"`
	AwemeCount []int64  `json:"aweme_count"`
}
