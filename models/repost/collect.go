package repost

import (
	"dongchamao/models/dcm"
	"dongchamao/models/entity"
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
	Undercarriage    int
	IsCoupon         int
}

type CollectAwemeRet struct {
	dcm.DcUserDyCollect
	AwemeAuthorID   string
	AwemeCover      string
	AwemeTitle      string
	AwemeCreateTime int64
	AwemeURL        string
	DiggCount       int64
	DiggCountIncr   int64
	AuthorAvatar    string
	AuthorNickname  string
}
type CollectShopRet struct {
	dcm.DcUserDyCollect
	Shop entity.DyShop `json:"shop"`
}

type CollectTagRet struct {
	dcm.DcUserDyCollectTag
	Count int64
}

type CollectCount struct {
	TagId int   `json:"tag_id" gorm:"tag_id"`
	Count int64 `json:"count" gorm:"count"`
}

type DyProductDailySlice struct {
	ProductId      string `json:"product_id"`
	DyProductDaily map[string]entity.DyProductDaily
}
