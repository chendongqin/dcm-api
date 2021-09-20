package entity

//
//var ShortVideoProductMap = HbaseEntity{
//	"image":          {String, "image"},
//	"saleroom":       {Double, "saleroom"},
//	"aweme_num":      {Long, "aweme_num"},
//	"cos_ratio":      {String, "cos_ratio"},
//	"price":          {Double, "price"},
//	"product_id":     {String, "product_id"},
//	"cos_fee":        {String, "cos_fee"},
//	"title":          {AJson, "title"},
//	"sales":          {Long, "sales"},
//	"platform_label": {AJson, "platform_label"},
//}
var ShortVideoCommodityTopNMap = HbaseEntity{
	"info":    {Long, "update_time"},
	"message": {AJson, "ranks"},
}

type ShortVideoCommodityTopN struct {
	UpdateTime int64               `json:"update_time"`
	Ranks      []ShortVideoProduct `json:"ranks"`
}
type ShortVideoProduct struct {
	Image         string `json:"image"`
	Saleroom      string `json:"saleroom"`
	AwemeNum      int64  `json:"aweme_num"`
	CosRatio      string `json:"cos_ratio"`
	Price         string `json:"price"`
	ProductId     string `json:"product_id"`
	CosFee        string `json:"cos_fee"`
	Title         string `json:"title"`
	Sales         int64  `json:"sales"`
	PlatformLabel string `json:"platform_label"`
}
