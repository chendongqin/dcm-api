package entity

var DyProductAwemeAuthorAnalysisMap = HbaseEntity{
	"product_id":     {String, "product_id"},
	"author_id":      {String, "author_id"},
	"nickname":       {String, "nickname"},
	"create_sdf":     {String, "create_sdf"},
	"display_id":     {String, "display_id"},
	"short_id":       {String, "short_id"},
	"score":          {Double, "score"},
	"level":          {Int, "level"},
	"first_name":     {String, "first_name"},
	"second_name":    {String, "second_name"},
	"avatar":         {String, "avatar"},
	"follow_count":   {Long, "follow_count"},
	"digg_count":     {Long, "digg_count"},
	"related_awemes": {AJson, "related_awemes"},
	"sales":          {Long, "sales"},
	"gmv":            {Double, "gmv"},
}

type DyProductAwemeAuthorAnalysis struct {
	ProductId     string                        `json:"product_id"`
	AuthorId      string                        `json:"author_id"`
	Nickname      string                        `json:"nickname"`
	CreateSdf     string                        `json:"create_sdf"`
	DisplayId     string                        `json:"display_id"`
	ShortId       string                        `json:"short_id"`
	Score         float64                       `json:"score"`
	Level         int                           `json:"level"`
	FirstName     string                        `json:"first_name"`
	SecondName    string                        `json:"second_name"`
	Avatar        string                        `json:"avatar"`
	FollowCount   int64                         `json:"follow_count"`
	RelatedAwemes []DyProductAuthorRelatedAweme `json:"related_awemes"`
	DiggCount     int64                         `json:"digg_count"`
	Sales         int64                         `json:"sales"`
	Gmv           float64                       `json:"gmv"`
}

type DyProductAuthorRelatedAweme struct {
	CommentCount    int64   `json:"comment_count"`
	AwemeTitle      string  `json:"aweme_title"`
	AwemeId         string  `json:"aweme_id"`
	AwemeUrl        string  `json:"aweme_url"`
	Sales           int64   `json:"sales"`
	AwemeGmv        float64 `json:"aweme_gmv"`
	DiggCount       int64   `json:"digg_count"`
	ForwardCount    int64   `json:"forward_count"`
	AwemeCover      string  `json:"aweme_cover"`
	AwemeCreateTime int64   `json:"aweme_create_time"`
}
