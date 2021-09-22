package es

type DyProductVideo struct {
	ProductId       string  `json:"product_id"`
	AwemeId         string  `json:"aweme_id"`
	AwemeTitle      string  `json:"aweme_title"`
	AwemeCover      string  `json:"aweme_cover"`
	AuthorId        string  `json:"author_id"`
	Avatar          string  `json:"avatar"`
	Nickname        string  `json:"nickname"`
	AwemeCreateTime int     `json:"aweme_create_time"`
	CommentCount    int     `json:"comment_count"`
	DiggCount       int     `json:"digg_count"`
	ForwardCount    int     `json:"forward_count"`
	Sales           int     `json:"sales"`
	AwemeGmv        float64 `json:"aweme_gmv"`
	AwemeCreateSdf  string  `json:"aweme_create_sdf"`
	AwemeUrl        string  `json:"aweme_url"`
}
