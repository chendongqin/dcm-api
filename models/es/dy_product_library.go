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
