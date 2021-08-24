package entity

const (
	String  = "string"
	Long    = "long"
	Float   = "float"
	Double  = "double"
	Byte    = "byte"
	Int     = "int"
	Bool    = "bool"
	MLong   = "m_long"
	MInt    = "m_int"
	MString = "m_string"
	MDouble = "m_double"
	Json    = "json"
	AJson   = "a_json"
)

type HbaseField struct {
	FieldType string
	FieldName string
}

type HbaseEntity map[string]HbaseField

var TestMap = HbaseEntity{
	"other_digg_count":         {Long, "digg_count"},
	"other_duration":           {Long, "duration"},
	"other_med_digg":           {Long, "med_digg"},
	"other_digg_follower_rate": {Double, "digg_follower_rate"},
	"other_artificial_data":    {Json, "artificial_data"},
}
