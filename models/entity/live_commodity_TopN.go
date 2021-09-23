package entity

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
