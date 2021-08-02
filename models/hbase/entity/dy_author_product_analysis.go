package entity

var DyAuthorProductAnalysisMap = HbaseEntity{
	"author_id":           {String, "author_id"},
	"product_id":          {String, "product_id"},
	"title":               {String, "title"},
	"avatar":              {String, "avatar"},
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
	"room_count":          {Int, "room_count"},
	"video_predict_gmv":   {Double, "video_predict_gmv"},
	"vedio_predict_sales": {Double, "vedio_predict_sales"},
	"vedio_count":         {Int, "vedio_count"},
	"row_time":            {String, "row_time"},
	"create_time":         {String, "create_time"},
	"shelf_time":          {Long, "shelf_time"},
	"vedio_product_sales": {AJson, "vedio_product_sales"},
	"room_product_sales":  {AJson, "room_product_sales"},
}

type DyAuthorProductAnalysis struct {
	AuthorId          string                        `json:"author_id"`
	ProductId         string                        `json:"product_id"`
	Title             string                        `json:"title"`
	Avatar            string                        `json:"avatar"`
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
	RoomCount         int                           `json:"room_count"`
	VideoPredictGmv   float64                       `json:"video_predict_gmv"`
	VedioCount        int                           `json:"vedio_count"`
	RowTime           string                        `json:"row_time"`
	CreateTime        string                        `json:"create_time"`
	ShelfTime         int64                         `json:"shelf_time"`
	Gmv               float64                       `json:"gmv"`
	Sales             float64                       `json:"sales"`
	ProductStatus     int                           `json:"product_status"`
	VedioPredictSales float64                       `json:"vedio_predict_sales"`
	VedioProductSales []DyProductAnalysisVideoSales `json:"vedio_product_sales"`
	RoomProductSales  []DyProductAnalysisRoomSales  `json:"room_product_sales"`
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
