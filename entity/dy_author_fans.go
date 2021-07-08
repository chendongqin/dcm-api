package entity

var DyAuthorFansMap = HbaseEntity{
	"follower_count":         {Long, "follower_count"},
	"total_fans_group_count": {Long, "total_fans_group_count"},
}

type DyAuthorFans struct {
	FollowerCount       int64 `json:"follower_count"`
	TotalFansGroupCount int64 `json:"total_fans_group_count"`
}
