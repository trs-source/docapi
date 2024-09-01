package docapi

type Schema struct {
	OneOf    []Ref     `json:"oneOf,omitempty"`
	Type     DataTypes `json:"type,omitempty"`
	Required []string  `json:"required,omitempty"`
	Items    *Items    `json:"items,omitempty"`
	// Preencher neste nível quando é object
	Properties Properties `json:"properties,omitempty"`
}

type Value map[string]string

type Items struct {
	OneOf []Ref     `json:"oneOf,omitempty"`
	Type  DataTypes `json:"type,omitempty"`
	Ref   string    `json:"$ref,omitempty"`
	// Preencher neste nível quando é array
	Properties Properties `json:"properties,omitempty"`
}

type Ref struct {
	Ref string `json:"$ref"`
}

// Chave primeiro map é nome do campo
// Chave segundo map -> type: valor
// Chave segundo map -> format: valor
// Chave segundo map -> $ref: valor. Usado quando é um sub model, slice ou não.
type Properties map[string]map[string]any

func (s *Schema) AddOneOfRef(modelName string, dataType DataTypes) {
	ref := Ref{"#/components/schemas/" + modelName}

	switch dataType {
	case DataTypeArray:
		if s.Items == nil {
			s.Items = &Items{Type: dataType}
		}
		s.Items.OneOf = append(s.Items.OneOf, ref)

	default:
		s.OneOf = append(s.OneOf, ref)
	}

}

func NewSchema(modelName string, dataType DataTypes) (schema *Schema) {
	schema = &Schema{}
	schema.AddOneOfRef(modelName, dataType)
	return
}
