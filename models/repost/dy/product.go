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
	IsCollect     int     `json:"is_collect"`
}

type ProductOrderDaily struct {
	Date       string  `json:"date"`
	OrderCount int64   `json:"order_count"`
	PvCount    int64   `json:"pv_count"`
	Rate       float64 `json:"rate"`
	AwemeNum   int     `json:"aweme_num"`
	RoomNum    int     `json:"room_num"`
	AuthorNum  int     `json:"author_num"`
}

type ProductOrderChart struct {
	Date       []string  `json:"date"`
	OrderCount []int64   `json:"order_count"`
	PvCount    []int64   `json:"pv_count"`
	Rate       []float64 `json:"rate"`
}

type ProductAuthorChart struct {
	Date             []string `json:"date"`
	AuthorCount      []int    `json:"author_count"`
	AwemeAuthorCount []int    `json:"aweme_author_count"`
	LiveAuthorCount  []int    `json:"live_author_count"`
}

type ProductLiveAwemeChart struct {
	Date       []string `json:"date"`
	LiveCount  []int    `json:"live_count"`
	AwemeCount []int    `json:"aweme_count"`
}

type DyProductLiveCount struct {
	Tags  []DyCate    `json:"tags"`
	Level []DyIntCate `json:"level"`
}

type ProductSalesTrendChart struct {
	DateTimestamp int64 `json:"date_timestamp"`
	Sales         int64 `json:"sales"`
	VideoNum      int   `json:"video_num"`
}

type DyProductAwemeCount struct {
	Tags  []DyCate    `json:"tags"`
	Level []DyIntCate `json:"level"`
}
