package entity

var DyProductAuthorAnalysisMap = HbaseEntity{
	"author_id":     {String, "author_id"},
	"display_id":    {String, "display_id"},
	"follow_count":  {Long, "follow_count"},
	"gmv":           {Double, "gmv"},
	"nickname":      {String, "nickname"},
	"price":         {Double, "price"},
	"product_id":    {String, "product_id"},
	"related_rooms": {AJson, "related_rooms"},
	"sales":         {Long, "sales"},
	"score":         {Double, "score"},
	"level":         {Int, "level"},
	"shop_tags":     {String, "shop_tags"},
	"short_id":      {String, "short_id"},
	"shop_id":       {String, "shop_id"},
}

type DyProductAuthorAnalysis struct {
	AuthorId     string                       `json:"author_id"`
	DisplayId    string                       `json:"display_id"`
	FollowCount  int64                        `json:"follow_count"`
	Gmv          float64                      `json:"gmv"`
	Nickname     string                       `json:"nickname"`
	Avatar       string                       `json:"avatar"`
	Price        float64                      `json:"price"`
	ProductId    string                       `json:"product_id"`
	RelatedRooms []DyProductAuthorRelatedRoom `json:"related_rooms"`
	Products     []DyAuthorProductDetail      `json:"products"`
	RoomNum      int                          `json:"room_num"`
	ProductNum   int                          `json:"product_num"`
	Sales        int64                        `json:"sales"`
	Score        float64                      `json:"score"`
	Level        int                          `json:"level"`
	ShopTags     string                       `json:"shop_tags"`
	ShortId      string                       `json:"short_id"`
	ShopId       string                       `json:"shop_id"`
	Date         string                       `json:"date"`
}

type DyAuthorProductDetail struct {
	Gmv       float64 `json:"gmv"`
	Price     float64 `json:"price"`
	ProductId string  `json:"product_id"`
	Sales     int64   `json:"sales"`
	Date      string  `json:"date"`
}

type DyProductAuthorRelatedRoom struct {
	EndTs      int64   `json:"end_ts"`
	Gmv        float64 `json:"gmv"`
	RoomId     string  `json:"room_id"`
	Sales      int64   `json:"sales"`
	StartTs    int64   `json:"start_ts"`
	Title      string  `json:"title"`
	Cover      string  `json:"cover"`
	TotalUser  int64   `json:"total_user"`
	LiveSecond int64   `json:"live_second"`
}
