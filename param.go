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

// https://swagger.io/docs/specification/serialization/
type Parameter struct {
	Required    bool    `json:"required,omitempty"`
	In          ParamIn `json:"in,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Example     string  `json:"example,omitempty"`
	ParamSchema *Schema `json:"schema,omitempty"`
}

type OptsParameter func(*Parameter)

func WithParamRequired() OptsParameter {
	return func(p *Parameter) {
		p.Required = true
	}
}

func WithParamDescription(description string) OptsParameter {
	return func(p *Parameter) {
		p.Description = description
	}
}

func WithParamExample(example string) OptsParameter {
	return func(p *Parameter) {
		p.Example = example
	}
}

// https://swagger.io/docs/specification/data-models/data-types/
type DataType string

const (
	schemaNone      DataType = ""
	DataTypeString  DataType = "string"
	DataTypeArray   DataType = "array"
	DataTypeBoolean DataType = "boolean"
	DataTypeInteger DataType = "integer"
	DataTypeNumber  DataType = "number"
	DataTypeObject  DataType = "object"
)

func (s DataType) String() string {
	return string(s)
}
