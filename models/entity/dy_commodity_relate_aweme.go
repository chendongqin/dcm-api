package entity

var DyCommodityRelateAwemeMap = HbaseEntity{
	"aweme_create_time": {String, "aweme_create_time"},
	"aweme_gmv":         {Float, "aweme_gmv"},
	"aweme_id":          {String, "aweme_id"},
	"aweme_title":       {String, "aweme_title"},
	"aweme_url":         {String, "aweme_url"},
	"aweme_cover":       {String, "aweme_cover"},
	"comment_count":     {Float, "comment_count"},
	"digg_count":        {Float, "digg_count"},
	"product_id":        {String, "product_id"},
	"sales":             {Float, "sales"},
	"share_countc":      {String, "share_countc"},
}

type DyCommodityRelateAweme struct {
	AwemeCreateTime string  `json:"aweme_create_time"`
	AwemeGmv        float64 `json:"aweme_gmv"`
	AwemeId         string  `json:"aweme_id"`
	AwemeTitle      string  `json:"aweme_title"`
	AwemeUrl        string  `json:"aweme_url"`
	AwemeCover      string  `json:"aweme_cover"`
	CommentCount    float64 `json:"comment_count"`
	DiggCount       float64 `json:"digg_count"`
	ProductId       string  `json:"product_id"`
	Sales           float64 `json:"sales"`
	ShareCountc     string  `json:"share_countc"`
}
