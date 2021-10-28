package es

type DyProduct struct {
	ProductId         string  `json:"product_id"`
	Image             string  `json:"image"`
	Title             string  `json:"title"`
	Price             float64 `json:"price"`
	CouponPrice       float64 `json:"coupon_price"`
	CommissionRate    float64 `json:"commission_rate"`
	Commission        float64 `json:"commission"`
	Pv                int64   `json:"pv"`
	OrderAccount      int64   `json:"order_account"` //昨日订单量
	Cvr               float64 `json:"cvr"`           //转化率
	WeekOrderAccount  int64   `json:"week_order_account"`
	MonthOrderAccount int64   `json:"month_order_account"`
	IsCoupon          int     `json:"is_coupon"`
	CommerceType      int     `json:"commerce_type"`
	PlatformLabel     string  `json:"platform_label"`
	ShopName          string  `json:"shop_name"`
	RelateAweme       int     `json:"relate_aweme"`
	RelateRoom        int     `json:"relate_room"`
	RelateAuthor      int     `json:"relate_author"`
	IsYesterday       int     `json:"is_yesterday"`
	DcmLevelFirst     string  `json:"dcm_level_first"`
	FirstCname        string  `json:"first_cname"`
	SecondCname       string  `json:"second_cname"`
	ThirdCname        string  `json:"third_cname"`
	IsStar            int     `json:"is_star"`
	IsCollect         int     `json:"is_collect"`
	Undercarriage     int     `json:"undercarriage"`
}

type ProductNew struct {
	ProductId       string `json:"product_id"`
	Image           string `json:"image"`
	Title           string `json:"title"`
	Price           string `json:"price"`
	CouponPrice     string `json:"coupon_price"`
	Commission      string `json:"commission"`
	DcmLevelFirst   string `json:"dcm_level_first"`
	FirstCname      string `json:"first_cname"`
	SecondCname     string `json:"second_cname"`
	ThirdCname      string `json:"third_cname"`
	Undercarriage   string `json:"undercarriage"`
	ShopId          string `json:"shop_id"`
	Pv              string `json:"pv"`
	Cvr             string `json:"cvr"`
	OrderAccount    string `json:"order_account"`
	Gpm             string `json:"gpm"`
	IsCoupon        string `json:"is_coupon"`
	PlatformLabel   string `json:"platform_label"`
	Pv7             string `json:"pv_7"`
	Cvr7            string `json:"cvr_7"`
	OrderAccount7   string `json:"order_account_7"`
	Gpm7            string `json:"gpm_7"`
	IsCoupon7       string `json:"is_coupon_7"`
	PlatformLabel7  string `json:"platform_label_7"`
	RelateAweme7    string `json:"relate_aweme_7"`
	RelateRoom7     string `json:"relate_room_7"`
	RelateAuthor7   string `json:"relate_author_7"`
	IsStar7         string `json:"is_star_7"`
	Pv15            string `json:"pv_15"`
	Cvr15           string `json:"cvr_15"`
	OrderAccount15  string `json:"order_account_15"`
	Gpm15           string `json:"gpm_15"`
	IsCoupon15      string `json:"is_coupon_15"`
	PlatformLabel15 string `json:"platform_label_15"`
	RelateAweme15   string `json:"relate_aweme_15"`
	RelateRoom15    string `json:"relate_room_15"`
	RelateAuthor15  string `json:"relate_author_15"`
	IsStar15        string `json:"is_star_15"`
	Dt              string `json:"dt"`
	OrderAccount30  string `json:"order_account_30"`
	RelateAuthor30  string `json:"relate_author_30"`
	IsStar30        string `json:"is_star_30"`
	PlatformLabel30 string `json:"platform_label_30"`
	Pv30            string `json:"pv_30"`
	RelateAweme30   string `json:"relate_aweme_30"`
	Gpm30           string `json:"gpm_30"`
	IsCoupon30      string `json:"is_coupon_30"`
	Cvr30           string `json:"cvr_30"`
	RelateRoom30    string `json:"relate_room_30"`
	CommerceType30  string `json:"commerce_type_30"`
	CommerceType15  string `json:"commerce_type_15"`
	CommerceType7   string `json:"commerce_type_7"`
	RelateAuthor    string `json:"relate_author"`
	RelateAweme     string `json:"relate_aweme"`
	CommerceType    string `json:"commerce_type"`
	IsStar          string `json:"is_star"`
	RelateRoom      string `json:"relate_room"`
	CommissionRate  string `json:"commission_rate"`
	IsCollect       int    `json:"is_collect"`
	ShopName        string `json:"shop_name"`
}
