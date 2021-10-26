package entity

var DyAuthorBasicMap = HbaseEntity{
	"follower_count":   {Long, "follower_count"},
	"total_fans_count": {Long, "total_fans_count"},
	"total_favorited":  {Long, "total_favorited"},
	"comment_count":    {Long, "comment_count"},
	"share_count":      {Long, "forward_count"},
}

type DyAuthorBasic struct {
	FollowerCount  int64 `json:"follower_count"`
	TotalFansCount int64 `json:"total_fans_count"`
	TotalFavorited int64 `json:"total_favorited"`
	CommentCount   int64 `json:"comment_count"`
	ForwardCount   int64 `json:"forward_count"`
}
