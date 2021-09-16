package entity

var DyAwemeCommentTopMap = HbaseEntity{
	"cid":         {String, "cid"},
	"text":        {String, "text"},
	"digg_count":  {Long, "digg_count"},
	"create_time": {String, "create_time"},
}

type DyAwemeCommentTop struct {
	Cid        string `json:"cid"`
	Text       string `json:"text"`
	DiggCount  int64  `json:"digg_count"`
	CreateTime string `json:"create_time"`
}
