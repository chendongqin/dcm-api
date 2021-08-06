package entity

var DyAuthorLiveTagsMap = HbaseEntity{
	"tags": {String, "tags"},
}

type DyAuthorLiveTags struct {
	Tags string `json:"tags"`
}
