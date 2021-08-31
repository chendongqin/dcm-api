package entity

var DyProductAwemeSalesTrendMap = HbaseEntity{
	"aweme_ids": {String, "aweme_ids"},
	"aweme_num": {Int, "aweme_num"},
	"sales":     {Long, "sales"},
}

type DyProductAwemeSalesTrend struct {
	AwemeIds string `json:"aweme_ids"`
	AwemeNum int    `json:"aweme_num"`
	Sales    int64  `json:"sales"`
}
