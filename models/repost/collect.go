package repost

import "dongchamao/models/dcm"

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
	ProductId string `json:"product_id"`
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
