package entity

var DyAuthorProductDateMappingMap = HbaseEntity{
	"product_list": {AJson, "product_list"},
}

type DyAuthorDateProductData struct {
	ProductList []DyAuthorDateProductList `json:"product_list"`
}

type DyAuthorDateProductList struct {
	PredictGmv   float64 `json:"predict_gmv"`
	PredictSales float64 `json:"predict_sales"`
	Price        float64 `json:"price"`
	ProductId    string  `json:"product_id"`
	ShelfTime    int64   `json:"shelf_time"`
}
