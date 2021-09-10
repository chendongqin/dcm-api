package dy

type AuthorVideoOverview struct {
	VideoNum        int64        `json:"video_num"`
	ProductVideo    int64        `json:"product_video"`
	AvgDiggCount    int64        `json:"avg_digg_count"`
	AvgCommentCount int64        `json:"avg_comment_count"`
	AvgForwardCount int64        `json:"avg_forward_count"`
	DurationChart   []VideoChart `json:"duration_chart"`
	PublishChart    []VideoChart `json:"publish_chart"`
	DiggChart       DateChart    `json:"digg_chart"`
	DiggMax         int64        `json:"digg_max"`
	DiggMin         int64        `json:"digg_min"`
	CommentChart    DateChart    `json:"comment_chart"`
	CommentMax      int64        `json:"comment_max"`
	CommentMin      int64        `json:"comment_min"`
	ForwardChart    DateChart    `json:"forward_chart"`
	ForwardMax      int64        `json:"forward_max"`
	ForwardMin      int64        `json:"forward_min"`
}

type VideoChart struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type DySimpleAweme struct {
	AuthorID        string `json:"author_id"`
	AwemeCover      string `json:"aweme_cover"`
	AwemeTitle      string `json:"aweme_title"`
	AwemeCreateTime int64  `json:"aweme_create_time"`
	AwemeURL        string `json:"aweme_url"`
	CommentCount    int64  `json:"comment_count"`
	DiggCount       int64  `json:"digg_count"`
	DownloadCount   int64  `json:"download_count"`
	Duration        int    `json:"duration"`
	ForwardCount    int64  `json:"forward_count"`
	ID              string `json:"id"`
	MusicID         string `json:"music_id"`
	ShareCount      int64  `json:"share_count"`
	PromotionNum    int    `json:"promotion_num"`
}

type DyAwemeProductSale struct {
	AwemeId       string  `json:"aweme_id"`
	ProductId     string  `json:"product_id"`
	Gmv           float64 `json:"gmv"`
	Sales         int64   `json:"sales"`
	Price         float64 `json:"price"`
	Title         string  `json:"title"`
	PlatformLabel string  `json:"platform_label"`
	ProductStatus int     `json:"product_status"`
	CouponInfo    string  `json:"coupon_info"`
	Image         string  `json:"image"`
}

type AuthorAwemeSum struct {
	Total      int     `json:"total"`
	Gmv        float64 `json:"gmv"`
	Sales      int64   `json:"sales"`
	AvgDigg    int64   `json:"avg_digg"`
	AvgShare   int64   `json:"avg_share"`
	AvgComment int64   `json:"avg_comment"`
}
