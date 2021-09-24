package entity

var DyAuthorStoreMap = HbaseEntity{
	"id":         {String, "id"},
	"sid":        {String, "sid"},
	"nick_name":  {String, "nick_name"},
	"shop_name":  {String, "shop_name"},
	"reputation": {Json, "reputation"},
	"brand_tag":  {String, "brand_tag"},
	"sec_uid":    {String, "sec_uid"},
	"brand":      {Int, "brand"},
}

type DyAuthorStore struct {
	Id         string            `json:"id"`
	Sid        string            `json:"sid"`
	NickName   string            `json:"nick_name"`
	ShopName   string            `json:"shop_name"`
	Reputation DyStoreReputation `json:"reputation"`
	BrandTag   string            `json:"brand_tag"`
	SecUid     string            `json:"sec_uid"`
	Brand      int               `json:"brand"`
}

type DyStoreReputation struct {
	Level   int     `json:"level"`
	Percent string  `json:"percent"`
	Sales   string  `json:"sales"`
	Score   float64 `json:"score"`
	Text    string  `json:"text"`
}
