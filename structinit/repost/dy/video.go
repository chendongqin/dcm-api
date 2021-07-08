package dy

type AuthorVideoOverview struct {
	VideoNum        int64        `json:"video_num"`
	ProductVideo    int64        `json:"product_video"`
	AvgDiggCount    int64        `json:"avg_digg_count"`
	AvgCommentCount int64        `json:"avg_comment_count"`
	AvgForwardCount int64        `json:"avg_forward_count"`
	DurationChart   []VideoChart `json:"duration_chart"`
	PublishChart    []VideoChart `json:"publish_chart"`
	DiggChart       []DateChart  `json:"digg_chart"`
	DiggMax         int64        `json:"digg_max"`
	DiggMin         int64        `json:"digg_min"`
	CommentChart    []DateChart  `json:"comment_chart"`
	CommentMax      int64        `json:"comment_max"`
	CommentMin      int64        `json:"comment_min"`
	ForwardChart    []DateChart  `json:"forward_chart"`
	ForwardMax      int64        `json:"forward_max"`
	ForwardMin      int64        `json:"forward_min"`
}

type VideoChart struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
