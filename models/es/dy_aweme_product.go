package es

type AwemeAuthorProductList struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
	Data     struct {
		Hits struct {
			Total    int         `json:"total"`
			MaxScore interface{} `json:"max_score"`
			Hits     []struct {
				Index  string             `json:"_index"`
				Type   string             `json:"_type"`
				Id     string             `json:"_id"`
				Score  interface{}        `json:"_score"`
				Source AwemeAuthorProduct `json:"_source"`
				Sort   []int              `json:"sort"`
			} `json:"hits"`
		} `json:"hits"`
	} `json:"data"`
	AwemeCreateTime struct {
		Value float64 `json:"value"`
	} `json:"aweme_create_time"`
	Sales struct {
		Value float64 `json:"value"`
	} `json:"sales"`
	AwemeGmv struct {
		Value float64 `json:"value"`
	} `json:"aweme_gmv"`
}
type AwemeAuthorProduct struct {
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
	AwemeCreateTime float64 `json:"aweme_create_time"`
	AwemeCreateSdf  string  `json:"aweme_create_sdf"`
	Nickname        string  `json:"nickname"`
	Avatar          string  `json:"avatar"`
	FollowerCount   int     `json:"follower_count"`
	ShortId         string  `json:"short_id"`
	UniqueId        string  `json:"unique_id"`
	AwemeCover      string  `json:"aweme_cover"`
	AwemeTitle      string  `json:"aweme_title"`
	DiggCount       int     `json:"digg_count"`
	CommentCount    int     `json:"comment_count"`
	ForwardCount    int     `json:"forward_count"`
	DownloadCount   int     `json:"download_count"`
	CrawlTime       int     `json:"crawl_time"`
	Duration        int     `json:"duration"`
	Score           float64 `json:"score"`
	Level           int     `json:"level"`
	Tags            string  `json:"tags"`
	TagsLevelTwo    string  `json:"tags_level_two"`
	Exist           int     `json:"exist"`
	Sales           float64 `json:"sales,omitempty"`
	AwemeGmv        float64 `json:"aweme_gmv,omitempty"`
	DistDate        string  `json:"dist_date,omitempty"`
}
