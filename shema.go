package docapi

import (
	"fmt"
	"strconv"
)

type Schema struct {
	OneOf    []Ref    `json:"oneOf,omitempty"`
	Required []string `json:"required,omitempty"`
	Type     DataType `json:"type,omitempty"`
	// Preencher neste nível quando é object
	Properties any    `json:"properties,omitempty"`
	Items      *Items `json:"items,omitempty"`
}

type Items struct {
	OneOf    []Ref    `json:"oneOf,omitempty"`
	Required []string `json:"required,omitempty"`
	Type     DataType `json:"type,omitempty"`
	// Preencher neste nível quando é array
	Properties any `json:"properties,omitempty"`
}

type Ref struct {
	Ref string `json:"$ref"`
}

type Property struct {
	Format   string   `json:"format,omitempty"`
	Type     DataType `json:"type,omitempty"`
	Items    *Items   `json:"items,omitempty"`
	Enum     []any    `json:"enum,omitempty"`
	Required []string `json:"required,omitempty"`
	Value    any      `json:"properties,omitempty"`
}

func (s *Schema) AddOneOfRef(modelName string, dataType DataType) {
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

// ConvertEnumType responsável por converter o valor do enum conforme dt.
func (p *Property) ConvertEnumType(dt DataType) {
	for i, e := range p.Enum {
		var (
			value any
			err   error
			s     = fmt.Sprint(e)
		)

		switch dt {
		case DataTypeInteger:
			value, err = strconv.Atoi(s)
		case DataTypeNumber:
			value, err = strconv.ParseFloat(s, 64)
		case DataTypeBoolean:
			value, err = strconv.ParseBool(s)
		default:
			value = e
		}

		if err != nil {
			value = e
		}

		p.Enum[i] = value
	}
}
