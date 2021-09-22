package entity

var ShortVideoCommodityTopNMap = HbaseEntity{
	"info":    {Long, "update_time"},
	"message": {AJson, "ranks"},
}

type ShortVideoCommodityTopN struct {
	UpdateTime int64               `json:"update_time"`
	Ranks      []ShortVideoProduct `json:"ranks"`
}
type ShortVideoProduct struct {
	Image         string  `json:"image"`
	Saleroom      float64 `json:"saleroom"`
	AwemeNum      int64   `json:"aweme_num"`
	CosRatio      float64 `json:"cos_ratio"`
	Price         float64 `json:"price"`
	ProductId     string  `json:"product_id"`
	CosFee        float64 `json:"cos_fee"`
	Title         string  `json:"title"`
	Sales         int64   `json:"sales"`
	PlatformLabel string  `json:"platform_label"`
}
