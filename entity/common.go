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
