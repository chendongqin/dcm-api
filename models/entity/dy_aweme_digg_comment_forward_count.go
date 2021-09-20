package entity

var DyAwemeDiggCommentForwardCountMap = HbaseEntity{
	"crawl_time":    {Long, "crawl_time"},
	"digg_count":    {Long, "digg_count"},
	"comment_count": {Long, "comment_count"},
	"share_count":   {Long, "forward_count"},
}

type DyAwemeDiggCommentForwardCount struct {
	CrawlTime    int64 `json:"crawl_time"`
	DiggCount    int64 `json:"digg_count"`
	CommentCount int64 `json:"comment_count"`
	ForwardCount int64 `json:"forward_count"`
}
