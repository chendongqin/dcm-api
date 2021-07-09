package entity

var DyAuthorBasicMap = HbaseEntity{
	"follower_count":          {Long, "follower_count"},
	"follower_count_before":   {Long, "follower_count_before"},
	"total_fans_count":        {Long, "total_fans_count"},
	"total_fans_count_before": {Long, "total_fans_count_before"},
	"total_favorited":         {Long, "total_favorited"},
	"total_favorited_before":  {Long, "total_favorited_before"},
	"comment_count":           {Long, "comment_count"},
	"comment_count_before":    {Long, "comment_count_before"},
	"forward_count":           {Long, "forward_count"},
	"forward_count_before":    {Long, "forward_count_before"},
}

type DyAuthorBasic struct {
	FollowerCount        int64 `json:"follower_count"`
	FollowerCountBefore  int64 `json:"follower_count_before"`
	TotalFansCount       int64 `json:"total_fans_count"`
	TotalFansCountBefore int64 `json:"total_fans_count_before"`
	TotalFavorited       int64 `json:"total_favorited"`
	TotalFavoritedBefore int64 `json:"total_favorited_before"`
	CommentCount         int64 `json:"comment_count"`
	CommentCountBefore   int64 `json:"comment_count_before"`
	ForwardCount         int64 `json:"forward_count"`
	ForwardCountBefore   int64 `json:"forward_count_before"`
}
