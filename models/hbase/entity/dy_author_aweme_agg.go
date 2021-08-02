package entity

var DyAuthorAwemeAggMap = HbaseEntity{
	"data": {AJson, "data"},
}

type DyAuthorAwemeAggData struct {
	Data []DyAuthorAwemeAgg `json:"data"`
}

type DyAuthorAwemeAgg struct {
	Duration        int64  `json:"duration"`
	CommentCount    int64  `json:"comment_count"`
	CrawlTime       int64  `json:"crawl_time"`
	AwemeID         string `json:"aweme_id"`
	DiggCount       int64  `json:"digg_count"`
	AwemeCreateTime int64  `json:"aweme_create_time"`
	ForwardCount    int64  `json:"forward_count"`
	DyPromotionId   string `json:"dy_promotion_id"`
}
