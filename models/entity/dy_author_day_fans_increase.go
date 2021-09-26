package entity

var DyAuthorDayFansIncreaseMap = HbaseEntity{
	"author_id":                {String, "author_id"},
	"short_id":                 {String, "short_id"},
	"unique_id":                {String, "unique_id"},
	"nickname":                 {String, "nickname"},
	"avatar":                   {String, "avatar"},
	"verification_type":        {String, "verification_type"},
	"verify_name":              {String, "verify_name"},
	"follow_count":             {String, "follower_count"},
	"author_id_date":           {String, "author_id_date"},
	"yesterday_follower_count": {String, "yesterday_follower_count"},
	"fans_inc":                 {String, "fans_inc"},
	"live_fans_inc":            {String, "live_fans_inc"},
	"aweme_fans_inc":           {String, "aweme_fans_inc"},
	"tags":                     {String, "tags"},
	"tags_level_two":           {String, "tags_level_two"},
	"rn_max":                   {String, "rn_max"},
}

type DyAuthorDayFansIncrease struct {
	AuthorId               string `json:"author_id"`
	ShortId                string `json:"short_id"`
	UniqueId               string `json:"unique_id"`
	Nickname               string `json:"nickname"`
	Avatar                 string `json:"avatar"`
	VerificationType       string `json:"verification_type"`
	VerifyName             string `json:"verify_name"`
	FollowerCount          string `json:"follower_count"`
	AuthorIdDate           string `json:"author_id_date"`
	YesterdayFollowerCount string `json:"yesterday_follower_count"`
	FansInc                string `json:"fans_inc"`
	LiveFansInc            string `json:"live_fans_inc"`
	AwemeFansInc           string `json:"aweme_fans_inc"`
	Tags                   string `json:"tags"`
	TagsLevelTwo           string `json:"tags_level_two"`
	RnMax                  string `json:"rn_max"`
}
