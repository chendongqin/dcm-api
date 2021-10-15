package entity

var DyAuthorRankMap = HbaseEntity{
	"fans_inc_date":        {String, "fans_inc_date"},
	"fans_inc_rn":          {String, "fans_inc_rn"},
	"predict_gmv_sum_date": {String, "predict_gmv_sum_date"},
	"predict_gmv_sum_rn":   {String, "predict_gmv_sum_rn"},
}

//达人排名
type DyAuthorRank struct {
	FansIncDate       string `json:"fans_inc_date"`
	FansIncRn         string `json:"fans_inc_rn"`
	PredictGmvSumDate string `json:"predict_gmv_sum_date"`
	PredictGmvSumRn   string `json:"predict_gmv_sum_rn"`
}
