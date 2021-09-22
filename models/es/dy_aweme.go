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
	DownloadCount   int64   `json:"download_count"`
	MusicId         string  `json:"music_id"`
	ShortId         string  `json:"short_id"`
	UniqueId        string  `json:"unique_id"`
	Nickname        string  `json:"nickname"`
	Exist           int64   `json:"exist"`
	AwemeGmv        float64 `json:"aweme_gmv"`
	Sales           int64   `json:"sales"`
	AwemeUrl        string  `json:"aweme_url"`
	ProductIds      string  `json:"product_ids"`
}
