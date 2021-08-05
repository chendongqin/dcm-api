package es

import "strings"

type EsDyLiveInfo struct {
	RoomId          string   `json:"room_id"`
	Nickname        string   `json:"nickname"`
	Avatar          string   `json:"avatar"`
	DisplayId       string   `json:"display_id"`
	ShortId         string   `json:"short_id"`
	Cover           string   `json:"cover"`
	Title           string   `json:"title"`
	Brand           int      `json:"brand"`
	RoomStatus      int      `json:"room_status"`
	AuthorId        string   `json:"author_id"`
	CreateTime      int64    `json:"create_time"`
	CreateTimestamp string   `json:"create_timestamp"`
	NumPromotions   int64    `json:"num_promotions"`
	PredictUvValue  float64  `json:"predict_uv_value"`
	RealUvValue     float64  `json:"real_uv_value"`
	PredictGmv      float64  `json:"predict_gmv"`
	RealGmv         float64  `json:"real_gmv"`
	AvgUserCount    float64  `json:"avg_user_count"`
	AllUserCount    float64  `json:"all_user_count"`
	Tags            string   `json:"tags"`
	TagsArr         []string `json:"tags_arr"`
}

func (receiver EsDyLiveInfo) GetTagsArr() []string {
	if receiver.Tags == "null" {
		receiver.Tags = ""
	}
	if receiver.Tags == "" {
		return []string{}
	}
	return strings.Split(receiver.Tags, "_")
}
