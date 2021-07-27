package dy

type SimpleDyProduct struct {
	ProductID     string  `json:"product_id"`
	Title         string  `json:"title"`
	MarketPrice   float64 `json:"market_price"`
	Price         float64 `json:"price"`
	URL           string  `json:"url"`
	Image         string  `json:"image"`
	Status        int     `json:"status"`
	ShopId        string  `json:"shop_id"`
	ShopName      string  `json:"shop_name"`
	Undercarriage int     `json:"undercarriage"`
	CrawlTime     int64   `json:"crawl_time"`
	PlatformLabel string  `json:"platform_label"`
	Label         string  `json:"label"`
	MinPrice      float64 `json:"min_price"`
	CosRatio      float64 `json:"cos_ratio"`
	CosRatioMoney float64 `json:"cos_ratio_money"`
}
