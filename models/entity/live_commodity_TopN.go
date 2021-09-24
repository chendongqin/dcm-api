package entity

var LiveCommodityTopNMap = HbaseEntity{
	"info":    {Long, "update_time"},
	"message": {AJson, "ranks"},
}

type LiveCommodityTopN struct {
	UpdateTime int64         `json:"update_time"`
	Ranks      []LiveProduct `json:"ranks"`
}

type DyLiveProductSaleTopRank struct {
	ProductId     string  `json:"product_id"`
	DateTime      string  `json:"date_time"`
	Sales         float64 `json:"sales"`
	Gmv           float64 `json:"gmv"`
	Title         string  `json:"title"`
	MarkerPrice   float64 `json:"marker_price"`
	Price         float64 `json:"price"`
	Images        string  `json:"images"`
	CosRatio      float64 `json:"cos_ratio"`
	CosFee        float64 `json:"cos_fee"`
	DcmCname      string  `json:"dcm_cname"`
	PlatformLabel string  `json:"platform_label"`
	Undercarriage int     `json:"undercarriage"`
	LiveCount     int64   `json:"live_count"`
}

type LiveProduct struct {
	Image         string  `json:"image"`
	Saleroom      float64 `json:"saleroom"`
	RoomNum       int64   `json:"room_num"`
	CosRatio      float64 `json:"cos_ratio"`
	Price         float64 `json:"price"`
	ProductId     string  `json:"product_id"`
	CosFee        float64 `json:"cos_fee"`
	Title         string  `json:"title"`
	Sales         int64   `json:"sales"`
	PlatformLabel string  `json:"platform_label"`
}
