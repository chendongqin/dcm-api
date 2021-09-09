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
	ProductId        string  `json:"product_id"`
	Image            string  `json:"image"`
	Price            float64 `json:"price"`
	CouponPrice      float64 `json:"coupon_price"`
	Pv               int64   `json:"pv"`
	OrderAccount     int64   `json:"order_account"` //昨日订单量
	WeekRelateAuthor int     `json:"week_relate_author"`
	WeekOrderAccount int64   `json:"week_order_account"`
	PlatformLabel    string  `json:"platform_label"`
}

type CollectAwemeRet struct {
	dcm.DcUserDyCollect
}

type CollectTagRet struct {
	dcm.DcUserDyCollectTag
	Count int64
}

type CollectCount struct {
	TagId int   `json:"tag_id" gorm:"tag_id"`
	Count int64 `json:"count" gorm:"count"`
}
