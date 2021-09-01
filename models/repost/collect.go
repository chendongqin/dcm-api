package repost

import "dongchamao/models/dcm"

type CollectRet struct {
	dcm.DcUserDyCollect
	FollowerCount      int64
	FollowerIncreCount int64
	Predict7Gmv        float64
	Predict7Digg       float64
	Avatar             string
}
