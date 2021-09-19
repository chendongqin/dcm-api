package entity

var DyProductCommentTopMap = HbaseEntity{
	"product_id": {String, "product_id"},
	"digg_info":  {String, "digg_info"},
}

type DyProductCommentTopStruct struct {
	ProductId string `json:"product_id"`
	DiggInfo  string `json:"digg_info"`
}

type DyProductCommentTop struct {
	DiggCount  string `json:"digg_count"`
	Cid        string `json:"cid"`
	Text       string `json:"text"`
	CreateTime string `json:"create_time"`
}
