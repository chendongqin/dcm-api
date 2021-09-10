package repost

import (
	"dongchamao/models/dcm"
)

type CollectAuthorRet struct {
	dcm.DcUserDyCollect
	FollowerCount      int64
	FollowerIncreCount int64
	Predict7Gmv        float64
	Predict7Digg       float64
	Avatar             string
}

type CollectProductRet struct {
	dcm.DcUserDyCollect
	ProductId        string
	Image            string
	Price            float64
	CouponPrice      float64
	Pv               int64
	OrderAccount     int64 //昨日订单量
	WeekRelateAuthor int
	WeekOrderAccount int64
	PlatformLabel    string
}

type CollectAwemeRet struct {
	dcm.DcUserDyCollect
	AwemeAuthorID   string `json:"author_id"`
	AwemeCover      string `json:"aweme_cover"`
	AwemeTitle      string `json:"aweme_title"`
	AwemeCreateTime int64  `json:"aweme_create_time"`
	AwemeURL        string `json:"aweme_url"`
	DiggCount       int64  `json:"digg_count"`
	AuthorAvatar    string `json:"avatar"`
	AuthorNickname  string `json:"nickname"`
}

type CollectTagRet struct {
	dcm.DcUserDyCollectTag
	Count int64
}

type CollectCount struct {
	TagId int   `json:"tag_id" gorm:"tag_id"`
	Count int64 `json:"count" gorm:"count"`
}
