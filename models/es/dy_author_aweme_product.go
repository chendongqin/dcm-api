package es

type EsDyAuthorAwemeProduct struct {
	AuthorId        string  `json:"author_id"`
	ProductId       string  `json:"product_id"`
	AwemeId         string  `json:"aweme_id"`
	Price           float64 `json:"price"`
	Sales           int     `json:"sales"`
	Gmv             float64 `json:"gmv"`
	PlatformLabel   string  `json:"platform_label"`
	Title           string  `json:"title"`
	Image           string  `json:"image"`
	BrandName       string  `json:"brand_name"`
	ShopName        string  `json:"shop_name"`
	ShopId          string  `json:"shop_id"`
	ShopIcon        string  `json:"shop_icon"`
	DcmLevelFirst   string  `json:"dcm_level_first"`
	FirstCname      string  `json:"first_cname"`
	SecondCname     string  `json:"second_cname"`
	ThirdCname      string  `json:"third_cname"`
	DistDate        string  `json:"dist_date"`
	AwemeCreateTime int     `json:"aweme_create_time"`
}
