package entity

var DyProductAwemeDailyDistributeMap = HbaseEntity{
	"aweme_id":   {String, "aweme_id"},
	"aweme_gmv":  {Double, "aweme_gmv"},
	"dist_date":  {String, "dist_date"},
	"price":      {Double, "price"},
	"product_id": {String, "product_id"},
	"sales":      {Long, "sales"},
}

type DyProductAwemeDailyDistribute struct {
	AwemeId   string  `json:"aweme_id"`
	AwemeGmv  float64 `json:"aweme_gmv"`
	DistDate  string  `json:"dist_date"`
	Price     float64 `json:"price"`
	ProductId string  `json:"product_id"`
	Sales     int64   `json:"sales"`
}
