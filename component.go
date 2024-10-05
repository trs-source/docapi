package docapi

import (
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Components struct {
	Examples Examples `json:"examples,omitempty"`
	// Chave: BearerAuth; BasicAuth; ApiKeyAuth; OAuth2;
	Security map[string]*SecuritySchemes `json:"securitySchemes,omitempty"`
	// A Chave é o nome do model/dto
	Schemas map[string]*Schema `json:"schemas,omitempty"`
}

// AddSecurity responsável por adicionar o ss no Components.
func (c *Components) AddSecurity(ss *SecuritySchemes) {
	if ss == nil {
		return
	}

	if len(c.Security) == 0 {
		c.Security = map[string]*SecuritySchemes{ss.TypeName: ss}
		return
	}

	c.Security[ss.TypeName] = ss
}

// AddSchemasAndExamples responsável por preencher components/schemas e content/contentType/shema, comforme modelo.
func (c *Components) AddSchemasAndExamples(modelValue reflect.Value, modelType reflect.Type, dataType DataType, opts ...OptsExample) (modelName string) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("[DocApi]", "Panic método AddSchemasAndExamples", err)
		}
	}()

	example := &Example{Summary: modelName}
	for _, fn := range opts {
		fn(example)
	}

	modelName = modelType.Name()
	if modelName == "" {
		modelName = example.TypeName
	}

	tokens, examples, properties, required := c.addSchemasAndExamples(modelValue, modelName, dataType, make(map[string]int, 0))
	example.Tokens = tokens
	example.Value = examples

	if len(c.Examples) == 0 {
		c.Examples = Examples{}
	}

	c.Examples[modelName] = example

	c.addSchema(modelName, &Schema{
		Type:       DataTypeObject,
		Properties: properties,
		Required:   required,
	})

	return
}

func (c *Components) addSchemasAndExamples(modValue reflect.Value, ownerName string, dataType DataType, navigation map[string]int) (tokens [][]byte, examples, properties any, required []string) {
	modelName := modValue.Type().Name()
	if modelName == "" {
		modelName = ownerName
	}

	// Essa validação evita loop infinito quando a struct tem auto relacionamento.
	seq, ok := navigation[modelName]
	if ok {
		seq += 1
		navigation[modelName] = seq
	} else {
		navigation[modelName] = 1
	}

	if seq > 2 {
		return
	}

	//examples
	examplesObject := make(map[string]any, 0)

	//schemas
	propValues := make(map[string]any, 0)

	for i := 0; i < modValue.NumField(); i++ {
		field := modValue.Type().Field(i)

		tagjson := field.Tag.Get("json")
		tagdocapi := field.Tag.Get("docapi")

		if tagjson == "-" {
			continue
		}

		if tagjson == "" {
			tagjson = field.Name
		}

		fieldType := field.Type

		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}

		// token usado para manter ordenado os valores que são adicionados em map.
		token := fmt.Sprintf("%d__%s$", i, uuid.New().String())
		tokens = append(tokens, []byte(token))

		var (
			newKindIsStruct, ok bool
			newTypeValue        reflect.Value
		)

		//Slice/Array
		if fieldType, newTypeValue, newKindIsStruct, ok = c.isSlice(fieldType); ok {
			property := &Property{Type: DataTypeArray, Items: &Items{}}
			propValues[token+tagjson] = property
			propertyType, exvalue, enum, isReq := c.parseFieldsAndTag(fieldType.Kind(), tagdocapi)

			if isReq {
				required = append(required, tagjson)
			}

			// Slice/Array de tipo primitivo
			if !newKindIsStruct {
				property.Items.Type = propertyType
				property.Enum = append(property.Enum, enum...)
				property.ConvertEnumType(propertyType)

				examplesObject[token+tagjson] = []any{exvalue}
				continue
			}

			tk, ex, prop, reqFields := c.addSchemasAndExamples(newTypeValue, tagjson, DataTypeArray, navigation)
			tokens = append(tokens, tk...)
			property.Items.Required = reqFields
			property.Items.Properties = prop

			examplesObject[token+tagjson] = ex
			continue
		}

		// Struct
		if newTypeValue, ok = c.isStruct(fieldType); ok {
			tk, ex, prop, reqFields := c.addSchemasAndExamples(newTypeValue, tagjson, DataTypeObject, navigation)
			tokens = append(tokens, tk...)

			propValues[token+tagjson] = &Property{
				Required: reqFields,
				Value:    prop,
			}

			examplesObject[token+tagjson] = ex
			continue
		}

		propertyType, exvalue, enum, isReq := c.parseFieldsAndTag(fieldType.Kind(), tagdocapi)

		if isReq {
			required = append(required, tagjson)
		}

		property := &Property{
			Type:   propertyType,
			Format: fieldType.Name(),
			Enum:   enum,
		}
		property.ConvertEnumType(propertyType)
		propValues[token+tagjson] = property

		examplesObject[token+tagjson] = exvalue
	}

	examples = examplesObject
	if dataType == DataTypeArray {
		examples = []map[string]any{examplesObject}
	}

	properties = propValues
	return
}

func (c *Components) addSchema(modelName string, schema *Schema) {
	if c.Schemas == nil {
		c.Schemas = make(map[string]*Schema, 1)
	}

	c.Schemas[modelName] = schema
}

// isSlice responsável por verificar se o type é slice ou array, caso seja, criar reflect.Value do tipo,
// também indica se o tipo é primitivo ([]int, []string...)
func (c *Components) isSlice(fieldType reflect.Type) (rt reflect.Type, newTypeValue reflect.Value, newKindIsStruct, ok bool) {
	rt = fieldType
	ok = rt.Kind() == reflect.Array || rt.Kind() == reflect.Slice
	if ok {
		rt = rt.Elem()
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}

		newTypeValue = reflect.ValueOf(reflect.New(rt).Interface()).Elem()

		newKindIsStruct = newTypeValue.Kind() == reflect.Struct
	}
	return
}

func (c *Components) isStruct(fieldType reflect.Type) (rt reflect.Value, ok bool) {
	ok = fieldType.Kind() == reflect.Struct && fieldType != reflect.TypeOf(time.Time{})
	if ok {
		rt = reflect.ValueOf(reflect.New(fieldType).Interface()).Elem()
	}
	return
}

// parseFieldsAndTag responsável por extrair os dados da tag docapi e o DataType conforme reflect.Kind.
func (c *Components) parseFieldsAndTag(fieldKind reflect.Kind, tagdocapi string) (pType DataType, exValue any, enum []any, required bool) {
	var example string

	for _, v := range strings.Split(tagdocapi, ";") {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "required:") {
			required = strings.TrimSpace(strings.TrimPrefix(v, "required:")) == "true"
			continue
		}

		if strings.HasPrefix(v, "example:") {
			example = strings.TrimSpace(strings.TrimPrefix(v, "example:"))
			continue
		}

		if strings.HasPrefix(v, "enum:") {
			for _, e := range strings.Split(strings.TrimPrefix(v, "enum:"), ",") {
				enum = append(enum, strings.TrimSpace(e))
			}
		}
	}

	var (
		defaultValue any
		err          error
	)

	switch fieldKind.String() {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		pType = DataTypeInteger
		defaultValue = 0
		exValue, err = strconv.Atoi(example)

	case "float32", "float64":
		pType = DataTypeNumber
		defaultValue = 0.1
		exValue, err = strconv.ParseFloat(example, 64)

	case "bool":
		pType = DataTypeBoolean
		defaultValue = false
		exValue, err = strconv.ParseBool(example)

	case "Time":
		pType = DataTypeString
		defaultValue = "datetime"
		exValue = example

	default:
		pType = DataTypeString
		defaultValue = "string"
		exValue = example
	}

	if err != nil || exValue == "" {
		exValue = defaultValue
	}

	return
}
