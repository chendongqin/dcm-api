package entity

var DyProductMap = HbaseEntity{
	"dy_promotion_id":        {String, "dy_promotion_id"},
	"title":                  {String, "title"},
	"market_price":           {Double, "market_price"},
	"price":                  {Double, "price"},
	"url":                    {String, "url"},
	"sales":                  {Long, "sales"},
	"image":                  {String, "image"},
	"count":                  {Long, "count"},
	"status":                 {Int, "status"},
	"shop_id":                {String, "shop_id"},
	"undercarriage":          {Int, "undercarriage"},
	"crawl_time":             {Long, "crawl_time"},
	"platform_label":         {String, "platform_label"},
	"coupon_end_time":        {String, "coupon_end_time"},
	"coupon_start_time":      {String, "coupon_start_time"},
	"tb_cat_leaf_name":       {String, "tb_cat_leaf_name"},
	"tb_cat_name":            {String, "tb_cat_name"},
	"tb_coupon_click_url":    {String, "tb_coupon_click_url"},
	"tb_coupon_info":         {String, "tb_coupon_info"},
	"tb_coupon_price":        {Double, "tb_coupon_price"},
	"tb_coupon_remain_count": {Long, "tb_coupon_remain_count"},
	"tb_h5_mprice":           {Double, "tb_h5_mprice"},
	"tb_h5_price":            {Double, "tb_h5_price"},
	"tb_item_url":            {String, "tb_item_url"},
	"tb_max_commission_rate": {Double, "tb_max_commission_rate"},
	"tb_nick":                {String, "tb_nick"},
	"tb_pic_url":             {String, "tb_pic_url"},
	"tb_sales":               {String, "tb_sales"},
	"tb_title":               {String, "tb_title"},
	"tb_user_type":           {Int, "tb_user_type"},
	"tb_volume":              {Long, "tb_volume"},
	"tb_zk_final_price":      {Double, "tb_zk_final_price"},
	"min_price":              {Double, "min_price"},
	"cos_ratio":              {Double, "cos_ratio"},
	"price_trends":           {AJson, "price_trends"},
	//"tb_small_images":        {AJson, "tb_small_images"},
}

type DyProduct struct {
	DyPromotionID       string                `json:"dy_promotion_id"`
	ProductID           string                `json:"product_id"`
	Title               string                `json:"title"`
	MarketPrice         float64               `json:"market_price"`
	Price               float64               `json:"price"`
	URL                 string                `json:"url"`
	Sales               int64                 `json:"sales"`
	Image               string                `json:"image"`
	Imgs                []string              `json:"imgs"`
	Count               int64                 `json:"count"`
	Status              int                   `json:"status"`
	ShopID              string                `json:"shop_id"`
	ShopName            string                `json:"shop_name"`
	Label               string                `json:"label"`
	Undercarriage       int                   `json:"undercarriage"`
	CrawlTime           int64                 `json:"crawl_time"`
	PlatformLabel       string                `json:"platform_label"`
	CouponEndTime       string                `json:"coupon_end_time"`
	CouponStartTime     string                `json:"coupon_start_time"`
	TbCatLeafName       string                `json:"tb_cat_leaf_name"`
	TbCatName           string                `json:"tb_cat_name"`
	TbCouponClickUrl    string                `json:"tb_coupon_click_url"`
	TbCouponInfo        string                `json:"tb_coupon_info"`
	TbCouponPrice       float64               `json:"tb_coupon_price"`
	TbCouponRemainCount int64                 `json:"tb_coupon_remain_count"`
	TbH5Mprice          float64               `json:"tb_h5_mprice"`
	TbH5Price           float64               `json:"tb_h5_price"`
	TbItemUrl           string                `json:"tb_item_url"`
	TbMaxCommissionRate float64               `json:"tb_max_commission_rate"`
	TbNick              string                `json:"tb_nick"`
	TbPicUrl            string                `json:"tb_pic_url"`
	TbSales             string                `json:"tb_sales"`
	TbTitle             string                `json:"tb_title"`
	TbUserType          int                   `json:"tb_user_type"`
	TbVolume            int64                 `json:"tb_volume"`
	TbZkFinalPrice      float64               `json:"tb_zk_final_price"`
	MinPrice            float64               `json:"min_price"`
	CosRatio            float64               `json:"cos_ratio"`
	PriceTrends         []DyProductPriceTrend `json:"price_trends"`
	//TbSmallImages       []interface{} `json:"tb_small_images"`
}

type DyProductPriceTrend struct {
	StartTime int64   `json:"start_time"`
	Price     float64 `json:"price"`
}