package entity

var DyAuthorDaySalesRankMap = HbaseEntity{
	"author_id":         {String, "author_id"},
	"short_id":          {String, "short_id"},
	"nickname":          {String, "nickname"},
	"avatar":            {String, "avatar"},
	"verification_type": {String, "verification_type"},
	"verify_name":       {String, "verify_name"},
	"predict_sales_sum": {String, "predict_sales_sum"},
	"predict_gmv_sum":   {String, "predict_gmv_sum"},
	"per_price":         {String, "per_price"},
	"room_id_count":     {String, "room_id_count"},
	"rn_max":            {String, "rn_max"},
	"tags":              {String, "tags"},
}

type DyAuthorDaySalesRank struct {
	AuthorId         string `json:"author_id"`
	ShortId          string `json:"short_id"`
	Nickname         string `json:"nickname"`
	Avatar           string `json:"avatar"`
	VerificationType string `json:"verification_type"`
	VerifyName       string `json:"verify_name"`
	PredictSalesSum  string `json:"predict_sales_sum"`
	PredictGmvSum    string `json:"predict_gmv_sum"`
	PerPrice         string `json:"per_price"`
	RoomIdCount      string `json:"room_id_count"`
	RnMax            string `json:"rn_max"`
	Tags             string `json:"tags"`
}
