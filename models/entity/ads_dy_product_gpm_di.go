package entity

var AdsDyProductGpmDiMap = HbaseEntity{
	"gpm": {Double, "gpm"},
}

type AdsDyProductGpmDi struct {
	Gpm float64 `json:"gpm"`
}
