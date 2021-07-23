package entity

var DyReputationMap = HbaseEntity{
	"uid":               {String, "uid"},
	"score":             {Double, "score"},
	"level":             {Int, "level"},
	"sales":             {String, "sales"},
	"last_30_day_sales": {Long, "last_30_day_sales"},
	"goods_score":       {Json, "goods_score"},
	"logistics_score":   {Json, "logistics_score"},
	"service_score":     {Json, "service_score"},
	"percentage":        {Double, "percentage"},
	"score_list":        {AJson, "score_list"},
	"dt_score_list":     {AJson, "dt_score_list"},
	"encrypt_shop_id":   {String, "encrypt_shop_id"},
	"shop_name":         {String, "shop_name"},
	"shop_logo":         {String, "shop_logo"},
	"blue":              {Int, "blue"},
	"crawl_time":        {Long, "crawl_time"},
}

type DyReputation struct {
	UID            string                       `json:"uid"`
	Score          float64                      `json:"score"`
	Level          int                          `json:"level"`
	Sales          string                       `json:"sales"`
	Last30DaySales int64                        `json:"last_30_day_sales"`
	GoodsScore     DyReputationScore            `json:"goods_score"`
	LogisticsScore DyReputationScore            `json:"logistics_score"`
	ServiceScore   DyReputationScore            `json:"service_score"`
	Percentage     float64                      `json:"percentage"`
	ScoreList      []DyReputationMonthScoreList `json:"score_list"`
	DtScoreList    []DyReputationDateScoreList  `json:"dt_score_list"`
	EncryptShopID  string                       `json:"encrypt_shop_id"`
	ShopName       string                       `json:"shop_name"`
	ShopLogo       string                       `json:"shop_logo"`
	Blue           int                          `json:"blue"`
	CrawlTime      int                          `json:"crawl_time"`
}

type DyReputationScore struct {
	Percentage float64 `json:"percentage"`
	Rating     string  `json:"rating"`
	Score      float64 `json:"score"`
}

type DyReputationMonthScoreList struct {
	Date  string  `json:"date"`
	Score float64 `json:"score"`
}

type DyReputationDateScoreList struct {
	Date    int     `json:"date"`
	DateStr string  `json:"date_str"`
	Score   float64 `json:"score"`
}
