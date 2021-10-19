package entity

var DyProductDealHxMap = HbaseEntity{
	//粉丝分析数据
	"gender":            {AJson, "gender"},
	"province":          {AJson, "province"},
	"city":              {AJson, "city"},
	"age_distrinbution": {AJson, "age_distrinbution"},
}

type DyProductDealHx struct {
	ProductID           string                   `json:"product_id"`
	//成交分析数据
	Gender           []DyAuthorFansGender   `json:"gender"`
	Province         []DyAuthorFansProvince `json:"province"`
	City             []DyAuthorFansCity     `json:"city"`
	AgeDistrinbution []DyAuthorFansAge      `json:"age_distrinbution"`
}
