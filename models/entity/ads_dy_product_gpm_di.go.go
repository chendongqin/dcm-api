package entity

var AdsDyProductPvDayDiMap = HbaseEntity{
	"pv": {Json, "pv"},
}

type AdsDyProductPvDayDi struct {
	Pv map[string]string `json:"pv"`
}
