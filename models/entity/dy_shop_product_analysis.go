package entity

var DyShopProductAnalysisMap = HbaseEntity{
	"product_id":      {String, "product_id"},
	"image":           {String, "image"},
	"title":           {String, "title"},
	"price":           {Double, "price"},
	"commission_rate": {Double, "commission_rate"},
	"gmv":             {Double, "gmv"},
	"sales":           {Long, "sales"},
	"month_pv_count":  {Long, "month_pv_count"},
	"month_cvr":       {Double, "month_cvr"},
	"first_cname":     {String, "first_cname"},
	"second_cname":    {String, "second_cname"},
	"dcm_level_first": {String, "dcm_level_first"},
}

type DyShopProductAnalysis struct {
	ProductId      string  `json:"product_id"`
	Image          string  `json:"image"`
	Title          string  `json:"title"`
	Price          float64 `json:"price"`
	CommissionRate float64 `json:"commission_rate"`
	Gmv            float64 `json:"gmv"`
	Sales          int64   `json:"sales"`
	ProductStatus  int `json:"product_status"`
	MonthPvCount   int64   `json:"month_pv_count"`
	MonthCvr       float64 `json:"month_cvr"`
	FirstCname     string  `json:"first_cname"`
	SecondCname    string  `json:"second_cname"`
	DcmLevelFirst  string  `json:"dcm_level_first"`
	Date           string  `json:"date"`
}
