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
	PredictUvValue  float64  `json:"predict_uv_value"`
	RealUvValue     float64  `json:"real_uv_value"`
	PredictGmv      float64  `json:"predict_gmv"`
	RealGmv         float64  `json:"real_gmv"`
	AvgUserCount    float64  `json:"avg_user_count"`
	AllUserCount    float64  `json:"all_user_count"`
	Tags            string   `json:"tags"`
	TagsArr         []string `json:"tags_arr"`
	NumProduct      int      `json:"num_product"`
	ProductTitle    string   `json:"product_title"`
	NumCrawlTimes   int      `json:"num_crawl_times"`
	PredictSales    float64  `json:"predict_sales"`
	RealSales       float64  `json:"real_sales"`
	DcmLevelFirst   string   `json:"dcm_level_first"`
	FirstCname      string   `json:"first_cname"`
	SecondCname     string   `json:"second_cname"`
	ThirdCname      string   `json:"third_cname"`
	MaxUserCount    int      `json:"max_user_count"`
	RowTime         string   `json:"row_time"`
	Dt              string   `json:"dt"`
	FinishTime      int64    `json:"finish_time"`
	WatchCnt        int64    `json:"watch_cnt"`
}

type NewEsDyLiveInfo struct {
	EsDyLiveInfo
	CustomerUnitPrice float64 `json:"customer_unit_price"`
	UserCount         int     `json:"user_count"`
	AvgStay           int     `json:"avg_stay"`
	AvgStayIndex      int     `json:"avg_stay_index"`
	FlowRatesIndex    string  `json:"flow_rates_index"`
	FlowRates         string  `json:"flow_rates"`
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

type EsDyLiveDetail struct {
	RoomId     string  `json:"room_id"`
	Nickname   string  `json:"nickname"`
	Avatar     string  `json:"avatar"`
	DisplayId  string  `json:"display_id"`
	ShortId    string  `json:"short_id"`
	Cover      string  `json:"cover"`
	Title      string  `json:"title"`
	RoomStatus int     `json:"room_status"`
	AuthorId   string  `json:"author_id"`
	CreateTime int64   `json:"create_time"`
	PredictGmv float64 `json:"predict_gmv"`
	FinishTime int64   `json:"finish_time"`
	WatchCnt   int64   `json:"watch_cnt"`
	AvgStay    float64 `json:"avg_stay"`
}
