package repost

type SearchData struct {
	SearchType string      `json:"search_type"`
	Note       string      `json:"note"`
	Content    interface{} `json:"content"`
}
