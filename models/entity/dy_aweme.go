package entity

var DyAwemeMap = HbaseEntity{
	"aweme_id":    {String, "aweme_id"},
	"crawl_time":  {Long, "crawl_time"},
	"aweme_title": {String, "aweme_title"},
	"data":        {Json, "data"},
}

type DyAweme struct {
	AwemeID    string      `json:"aweme_id"`
	CrawlTime  int         `json:"crawl_time"`
	AwemeTitle string      `json:"aweme_title"`
	Data       DyAwemeData `json:"data"`
}

type DyAwemeData struct {
	AuthorID        string         `json:"author_id"`
	AwemeCover      string         `json:"aweme_cover"`
	AwemeTitle      string         `json:"aweme_title"`
	AwemeCreateTime int64          `json:"aweme_create_time"`
	AwemeURL        string         `json:"aweme_url"`
	CommentCount    int64          `json:"comment_count"`
	DiggCount       int64          `json:"digg_count"`
	DownloadCount   int64          `json:"download_count"`
	Duration        int            `json:"duration"`
	DyPromotionID   []string       `json:"dy_promotion_id"`
	ForwardCount    int64          `json:"forward_count"`
	Sales           []int64        `json:"sales"`
	ID              string         `json:"id"`
	MusicID         string         `json:"music_id"`
	PromotionID     []string       `json:"promotion_id"`
	ShareCount      int64          `json:"share_count"`
	Topic           []DyAwemeTopic `json:"topic"`
}

type DyAwemeTopic struct {
	IsCommerce int    `json:"is_commerce"`
	TagID      string `json:"tag_id"`
	TagName    string `json:"tag_name"`
}
