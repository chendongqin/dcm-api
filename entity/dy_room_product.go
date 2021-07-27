package entity

var DyRoomProductMap = HbaseEntity{
	"trend_data":                       {AJson, "trend_data"},
	"price":                            {Double, "price"},
	"author_id":                        {String, "author_id"},
	"room_id":                          {String, "room_id"},
	"other_predict_sales_trend":        {AJson, "predict_sales_trend"},
	"other_predict_sales_detail_trend": {AJson, "predict_sales_detail_trend"},
}

type DyRoomProduct struct {
	TrendData        []DyRoomProductTrend           `json:"trend_data"`
	SalesTrend       []DyRoomProductSaleTrend       `json:"sales_trend"`
	SalesDetailTrend []DyRoomProductSaleDetailTrend `json:"sales_detail_trend"`
	Price            float64                        `json:"price"`
	AuthorId         string                         `json:"author_id"`
	RoomId           string                         `json:"room_id"`
}

type DyRoomProductSaleTrend struct {
	PredictSales float64 `json:"predict_sales"`
	EndTime      int64   `json:"endTime"`
	PredictGmv   float64 `json:"predict_gmv"`
}

type DyRoomProductSaleDetailTrend struct {
	PredictSales float64 `json:"predict_sales"`
	StartTime    int64   `json:"start_time"`
	EndTime      int64   `json:"endTime"`
	PredictGmv   float64 `json:"predict_gmv"`
}

type DyRoomProductTrend struct {
	CrawlTime int64   `json:"crawl_time"`
	Price     float64 `json:"price"`
	Sales     float64 `json:"sales"`
}
