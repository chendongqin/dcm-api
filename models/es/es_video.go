package es

type DyProductVideo struct {
	AuthorId        string  `json:"author_id"`
	ProductId       string  `json:"product_id"`
	AwemeId         string  `json:"aweme_id"`
	ShopId          string  `json:"shop_id"`
	ShopName        string  `json:"shop_name"`
	ShopIcon        string  `json:"shop_icon"`
	Price           float64 `json:"price"`
	PlatformLabel   string  `json:"platform_label"`
	Title           string  `json:"title"`
	Image           string  `json:"image"`
	BrandName       string  `json:"brand_name"`
	DcmLevelFirst   string  `json:"dcm_level_first"`
	FirstCname      string  `json:"first_cname"`
	SecondCname     string  `json:"second_cname"`
	ThirdCname      string  `json:"third_cname"`
	AwemeCreateTime int64   `json:"aweme_create_time"`
	AwemeCreateSdf  string  `json:"aweme_create_sdf"`
	Nickname        string  `json:"nickname"`
	Avatar          string  `json:"avatar"`
	FollowerCount   int64   `json:"follower_count"`
	ShortId         string  `json:"short_id"`
	UniqueId        string  `json:"unique_id"`
	AwemeCover      string  `json:"aweme_cover"`
	AwemeTitle      string  `json:"aweme_title"`
	DiggCount       int64   `json:"digg_count"`
	CommentCount    int64   `json:"comment_count"`
	ForwardCount    int64   `json:"forward_count"`
	DownloadCount   int64   `json:"download_count"`
	CrawlTime       int64   `json:"crawl_time"`
	Duration        int64   `json:"duration"`
	Score           float64 `json:"score"`
	Level           int     `json:"level"`
	Tags            string  `json:"tags"`
	TagsLevelTwo    string  `json:"tags_level_two"`
	Exist           int     `json:"exist"`
	AwemeGmv        float64 `json:"aweme_gmv"`
	DistDate        string  `json:"dist_date"`
	Sales           int64   `json:"sales"`
	AwemeUrl        string  `json:"aweme_url"`
}

type SumProductVideoAuthor struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
	Data     struct {
		Hits struct {
			Total int `json:"total"`
			Hits  []struct {
				Source DyProductVideo `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	} `json:"data"`
	Sales struct {
		Value int64 `json:"value"`
	} `json:"sales"`
	AwemeGmv struct {
		Value float64 `json:"value"`
	} `json:"aweme_gmv"`
}

type CountAuthorProductAweme struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
	Awemes   struct {
		Buckets []interface{} `json:"buckets"`
	} `json:"awemes"`
	Products struct {
		Buckets []interface{} `json:"buckets"`
	} `json:"products"`
}
