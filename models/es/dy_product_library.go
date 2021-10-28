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
	ProductId      string  `json:"product_id"`
	Image          string  `json:"image"`
	Title          string  `json:"title"`
	Price          float64 `json:"price"`
	CouponPrice    float64 `json:"coupon_price"`
	Commission     float64 `json:"commission"`
	DcmLevelFirst  string  `json:"dcm_level_first"`
	FirstCname     string  `json:"first_cname"`
	SecondCname    string  `json:"second_cname"`
	ThirdCname     string  `json:"third_cname"`
	Undercarriage  int64   `json:"undercarriage"`
	ShopId         string  `json:"shop_id"`
	CommissionRate float64 `json:"commission_rate"`
	Pv             int64   `json:"pv"`
	Cvr            float64 `json:"cvr"`
	OrderAccount   int64   `json:"order_account"`
	Gpm            int64   `json:"gpm"`
	IsCoupon       int     `json:"is_coupon"`
	CommerceType   int64   `json:"commerce_type"`
	PlatformLabel  string  `json:"platform_label"`
	RelateAweme    int64   `json:"relate_aweme"`
	RelateRoom     int64   `json:"relate_room"`
	RelateAuthor   int64   `json:"relate_author"`
	IsStar         int     `json:"is_star"`
	OrderAccount7  int64   `json:"order_account_7"`
	RelateRoom7    int64   `json:"relate_room_7"`
	Pv7            int64   `json:"pv_7"`
	CommerceType7  int64   `json:"commerce_type_7"`
	Gpm7           int64   `json:"gpm_7"`
	RelateAweme7   int64   `json:"relate_aweme_7"`
	Cvr7           float64 `json:"cvr_7"`
	RelateAuthor7  int64   `json:"relate_author_7"`
	ShopName       string  `json:"shop_name"`
	IsCollect      int     `json:"is_collect"`
}
