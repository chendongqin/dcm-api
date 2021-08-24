package es

type DyAuthorFollowerTop struct {
	AuthorCover           string `json:"author_cover"`
	AuthorID              string `json:"author_id"`
	AwemeIncFollowerCount int    `json:"aweme_inc_follower_count"`
	City                  string `json:"city"`
	Country               string `json:"country"`
	DateTime              string `json:"date_time"`
	FollowerCount         int    `json:"follower_count"`
	Gender                int    `json:"gender"`
	ID                    string `json:"id"`
	IncFollowerCount      int    `json:"inc_follower_count"`
	IsDelivery            int    `json:"is_delivery"`
	LiveIncFollowerCount  int    `json:"live_inc_follower_count"`
	Nickname              string `json:"nickname"`
	Province              string `json:"province"`
	ShortID               string `json:"short_id"`
	Tags                  string `json:"tags"`
	TagsLevelTwo          string `json:"tags_level_two"`
	UniqueID              string `json:"unique_id"`
	VerificationType      int    `json:"verification_type"`
	VerifyName            string `json:"verify_name"`
}
