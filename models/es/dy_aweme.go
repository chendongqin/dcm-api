package es

type DyAweme struct {
	AwemeId         string  `json:"aweme_id"`
	AwemeTitle      string  `json:"aweme_title"`
	AwemeCover      string  `json:"aweme_cover"`
	AwemeCreateTime int64   `json:"aweme_create_time"`
	DiggCount       int64   `json:"digg_count"`
	CommentCount    int64   `json:"comment_count"`
	ShareCount      int64   `json:"share_count"`
	CrawlTime       int64   `json:"crawl_time"`
	Duration        int64   `json:"duration"`
	AuthorId        string  `json:"author_id"`
	Avatar          string  `json:"avatar"`
	DistDate        string  `json:"dist_date"`
	AwemeGmv        float64 `json:"aweme_gmv"`
	Sales           int64   `json:"sales"`
}
