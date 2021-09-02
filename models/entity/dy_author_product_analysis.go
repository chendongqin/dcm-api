package entity

var DyAuthorProductAnalysisMap = HbaseEntity{
	"author_id":           {String, "author_id"},
	"product_id":          {String, "product_id"},
	"title":               {String, "title"},
	"image":               {String, "image"},
	"price":               {Double, "price"},
	"shop_id":             {String, "shop_id"},
	"shop_name":           {String, "shop_name"},
	"shop_icon":           {String, "shop_icon"},
	"brand_name":          {String, "brand_name"},
	"platform":            {String, "platform"},
	"dcm_level_first":     {String, "dcm_level_first"},
	"first_cname":         {String, "first_cname"},
	"second_cname":        {String, "second_cname"},
	"third_cname":         {String, "third_cname"},
	"live_predict_sales":  {Double, "live_predict_sales"},
	"live_predict_gmv":    {Double, "live_predict_gmv"},
	"room_count":          {Long, "room_count"},
	"aweme_predict_gmv":   {Double, "aweme_predict_gmv"},
	"aweme_predict_sales": {Double, "aweme_predict_sales"},
	"aweme_count":         {Long, "aweme_count"},
	"shelf_time":          {Long, "shelf_time"},
	"status":              {Int, "status"},
	"aweme_list":          {AJson, "aweme_list"},
	"room_list":           {AJson, "room_list"},
}

type DyAuthorProductAnalysis struct {
	AuthorId          string                        `json:"author_id"`
	ProductId         string                        `json:"product_id"`
	Title             string                        `json:"title"`
	Image             string                        `json:"image"`
	Price             float64                       `json:"price"`
	ShopId            string                        `json:"shop_id"`
	ShopName          string                        `json:"shop_name"`
	ShopIcon          string                        `json:"shop_icon"`
	BrandName         string                        `json:"brand_name"`
	Platform          string                        `json:"platform"`
	DcmLevelFirst     string                        `json:"dcm_level_first"`
	FirstCname        string                        `json:"first_cname"`
	SecondCname       string                        `json:"second_cname"`
	ThirdCname        string                        `json:"third_cname"`
	LivePredictSales  float64                       `json:"live_predict_sales"`
	LivePredictGmv    float64                       `json:"live_predict_gmv"`
	RoomCount         int64                         `json:"room_count"`
	AwemePredictGmv   float64                       `json:"aweme_predict_gmv"`
	AwemeCount        int64                         `json:"aweme_count"`
	ShelfTime         int64                         `json:"shelf_time"`
	Gmv               float64                       `json:"gmv"`
	Sales             float64                       `json:"sales"`
	Status            int                           `json:"status"`
	AwemePredictSales float64                       `json:"aweme_predict_sales"`
	AwemeList         []DyProductAnalysisVideoSales `json:"aweme_list"`
	RoomList          []DyProductAnalysisRoomSales  `json:"room_list"`
}

type DyProductAnalysisRoomSales struct {
	RoomProductId string  `json:"room_product_id"`
	PredictSales  float64 `json:"predict_sales"`
	PredictGmv    float64 `json:"predict_gmv"`
}

type DyProductAnalysisVideoSales struct {
	RoomProductId string  `json:"room_product_id"`
	PredictSales  float64 `json:"predict_sales"`
	PredictGmv    float64 `json:"predict_gmv"`
}
