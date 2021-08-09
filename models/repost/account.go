package repost

type SearchData struct {
	Id         int         `json:"id"`
	SearchType string      `json:"search_type"`
	Note       string      `json:"note"`
	Content    interface{} `json:"content"`
}
