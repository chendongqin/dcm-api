package entity

var DyAwemeCommentTopMap = HbaseEntity{
	"aweme_id":  {String, "aweme_id"},
	"digg_info": {String, "digg_info"},
}

type DyAwemeCommentTopStruct struct {
	AwemeId  string `json:"aweme_id"`
	DiggInfo string `json:"digg_info"`
}

type DyAwemeCommentTop struct {
	DiggCount  string `json:"digg_count"`
	Cid        string `json:"cid"`
	Text       string `json:"text"`
	CreateTime string `json:"create_time"`
	TagName    string `json:"tag_name"`
}
