package entity

var DyLivePmtMap = HbaseEntity{
	"room_status":  {Int, "room_status"},
	"author_id":    {String, "author_id"},
	"room_id":      {String, "room_id"},
	"create_time":  {Long, "create_time"},
	"crawl_time":   {Long, "crawl_time"},
	"purchase_cnt": {Long, "purchase_cnt"},
	"cur":          {String, "cur"},
	"promotions":   {AJson, "promotions"},
	"top":          {Int, "top"},
	"is_bubble":    {Bool, "is_bubble"},
}

type DyLivePmt struct {
	RoomStatus  int               `json:"room_status"`
	AuthorID    string            `json:"author_id"`
	RoomID      string            `json:"room_id"`
	CreateTime  int64             `json:"create_time"`
	CrawlTime   int64             `json:"crawl_time"`
	PurchaseCnt int               `json:"purchase_cnt"`
	Cur         string            `json:"cur"`
	Promotions  []DyLivePromotion `json:"promotions"`
	Top         int               `json:"top"`
	IsBubble    bool              `json:"is_bubble"`
}

type DyLivePromotion struct {
	DyPromotionID  string  `json:"dy_promotion_id"`  //抖音商品id
	ProductID      string  `json:"product_id"`       //第三方商品id
	ForSale        int     `json:"for_sale"`         //商品状态 0:刚上架 2:在售 4:下架
	StartTime      int64   `json:"start_time"`       //上架时间
	StopTime       int64   `json:"stop_time"`        //下架时间
	StartUserCount int64   `json:"start_user_count"` //上架uv
	StartTotalUser int64   `json:"start_total_user"` //上架总pv
	EndUserCount   int64   `json:"end_user_count"`   //结束uv
	EndTotalUser   int64   `json:"end_total_user"`   //结束总pv
	Price          float64 `json:"price"`            //价格
	Coupon         float64 `json:"coupon"`           //优惠券
	CouponHeader   string  `json:"coupon_header"`
	InitialSales   int64   `json:"initial_sales"` //初始销量
	FinalSales     int64   `json:"final_sales"`   //结束销量
	Sales          int64   `json:"sales"`         //当前销量
	InStock        bool    `json:"in_stock"`      //是否有库存
	Title          string  `json:"title"`
	ElasticTitle   string  `json:"elastic_title"` //优惠信息
	CosRatio       float64 `json:"cos_ratio"`
	Source         string  `json:"source"`   //来源
	ExtInfo        string  `json:"ext_info"` //优惠券信息
	Cover          string  `json:"cover"`    //封面
	InitialPv      int64   `json:"initial_pv"`
	FinalPv        int64   `json:"final_pv"`
	Pv             int64   `json:"pv"`
	Campaign       bool    `json:"campaign"`
	Index          int     `json:"index"` //第几个商品
	ShopID         string  `json:"shop_id"`
	FlushBuy       bool    `json:"flush_buy"`
	BubbleDuration int     `json:"bubble_duration"`
	BubblePv       int     `json:"bubble_pv"`
	HasH5PmtInfo   bool    `json:"has_h5_pmt_info"`
}
