package entity

var DyCommodityRelateLiveMap = HbaseEntity{
	"create_time": {Long, "create_time"},
	"gmv":         {Float, "gmv"},
	"room_id":     {String, "room_id"},
	"room_title":  {String, "title"},
	"room_cover":  {String, "cover"},
	"total_user":  {Long, "max_user_count"},
	"product_id":  {String, "product_id"},
	"sales":       {Long, "sales"},
	"price":       {Float, "price"},
}

type DyCommodityRelateLive struct {
	RoomId       string  `json:"room_id"`
	ProductId    string  `json:"product_id"`
	Cover        string  `json:"cover"`
	CreateTime   int64   `json:"create_time"`
	Title        string  `json:"title"`
	MaxUserCount int64   `json:"max_user_count"`
	Gmv          float64 `json:"gmv"`
	Sales        int64   `json:"sales"`
}
