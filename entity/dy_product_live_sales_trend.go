package entity

var DyProductLiveSalesTrendMap = HbaseEntity{
	"price":    {Double, "price"},
	"room_num": {Int, "room_num"},
	"sales":    {Long, "sales"},
	"room_ids": {String, "room_ids"},
}

type DyProductLiveSalesTrend struct {
	Price   float64 `json:"price"`
	RoomNum int     `json:"room_num"`
	Sales   int64   `json:"sales"`
	RoomIds string  `json:"room_ids"`
}
