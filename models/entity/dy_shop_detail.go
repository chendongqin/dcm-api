package entity

var DyShopDetailMap = HbaseEntity{
	"shop_id":         {String, "shop_id"},
	"sales":           {Long, "sales"},
	"gmv":             {Double, "gmv"},
	"classifications": {Json, "classifications"},
	"price_dist":      {Json, "price_dist"},
	"aweme_num":       {Long, "aweme_num"},
	"live_num":        {Long, "live_num"},
	"30d_aweme_cnt":   {Long, "30d_aweme_cnt"},
	"30d_live_cnt":    {Long, "30d_live_cnt"},
	"30d_author_cnt":  {Long, "30d_author_cnt"},
	"product_cnt":     {Long, "product_cnt"},
	"shop_cname":      {String, "shop_cname"},
	"30d_sales":       {Long, "30d_sales"},
	"30d_gmv":         {Double, "30d_gmv"},
	"30d_pct":         {Double, "30d_pct"},
}

type DyShopDetail struct {
	ShopId          string         `json:"shop_id"`
	Sales           int64          `json:"sales"`
	Gmv             float64        `json:"gmv"`
	Classifications map[string]int `json:"classifications"`
	PriceDist       map[string]int `json:"price_dist"`
	AwemeNum        int64          `json:"aweme_num"`
	LiveNum         int64          `json:"live_num"`
	D30AwemeCnt     int64          `json:"30d_aweme_cnt"`
	D30LiveCnt      int64          `json:"30d_live_cnt"`
	D30AuthorCnt    int64          `json:"30d_author_cnt"`
	ShopName        string         `json:"shop_name"`
	D30Sales        int64          `json:"30d_sales"`
	D30Gmv          int64          `json:"30d_gmv"`
	D30Pct          int64          `json:"30d_pct"`
}
