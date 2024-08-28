package docapi

// https://swagger.io/docs/specification/serialization/
type ParamIn string

const (
	ParamQuery  ParamIn = "query"
	ParamPath   ParamIn = "path"
	ParamHeader ParamIn = "header"
	ParamCookie ParamIn = "cookie"
)

func (p ParamIn) String() string {
	return string(p)
}

// https://swagger.io/docs/specification/data-models/data-types/
type DataTypes string

const (
	schemaNone      DataTypes = ""
	DataTypeString  DataTypes = "string"
	DataTypeArray   DataTypes = "array"
	DataTypeBoolean DataTypes = "boolean"
	DataTypeInteger DataTypes = "integer"
	DataTypeNumber  DataTypes = "number"
	DataTypeObject  DataTypes = "object"
)

func (s DataTypes) String() string {
	return string(s)
}

func GetDataType(fieldTpName string) (fieldType, fieldFormat string, value any) {
	fieldFormat = fieldTpName
	switch fieldTpName {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		fieldType = DataTypeInteger.String()
		value = 0

	case "float32", "float64":
		fieldType = DataTypeNumber.String()
		value = 0.1

	case "bool":
		fieldType = DataTypeBoolean.String()
		value = true

	default:
		fieldType = DataTypeString.String()
		value = "string"
	}
	return
}
