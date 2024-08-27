// Versão do Swagger 3.0.
//
// https://swagger.io/docs/specification/about/

package openapi

import (
	"reflect"
	"strings"
	"time"
)

// A chave deve ser a url do swagger concatenado com doc.json.
//
// Exemplo: http://localhost:8080/context-path/swagger/doc.json
type Docs map[string]*Doc

type Doc struct {
	Key          string        `json:"-"`
	Version      string        `json:"openapi"`
	Servers      []Servers     `json:"servers,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
	Info         *Info         `json:"info,omitempty"`
	// a chave é o endpoint
	Paths      map[string]Path `json:"paths,omitempty"`
	Components *Components     `json:"components,omitempty"`
}

type Servers struct {
	Url string `json:"url"`
}

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
}

type Info struct {
	Description string   `json:"description,omitempty"`
	Version     string   `json:"version,omitempty"`
	Title       string   `json:"title,omitempty"`
	Contact     *Contact `json:"contact,omitempty"`
	License     *License `json:"license,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type Components struct {
	Examples Examples `json:"examples,omitempty"`
	// Chave: BearerAuth; BasicAuth; ApiKeyAuth; OAuth2;
	Security map[string]*SecuritySchemes `json:"securitySchemes,omitempty"`
	// A Chave é o nome do model/dto
	Schemas map[string]*Schema `json:"schemas,omitempty"`
}

func (d *Doc) AddServer(url ...string) {
	for _, v := range url {
		d.Servers = append(d.Servers, Servers{v})
	}
}

func (s *Doc) AddPath(method, pattern string, pathStructure *PathsStructure) {
	method = strings.ToLower(method)
	// A raiz do path é a url e dentro contém os métodos get, post, put...
	// Se localizar a url, então adiciona o método.
	if len(s.Paths) == 0 {
		s.Paths = map[string]Path{pattern: {method: pathStructure}}
		return
	}

	if paths, ok := s.Paths[pattern]; ok {
		paths[method] = pathStructure
		s.Paths[pattern] = paths

	} else {
		s.Paths[pattern] = Path{method: pathStructure}
	}
}

func (d *Doc) AddComponentSecurity(ss *SecuritySchemes) {
	if ss == nil {
		return
	}

	if len(d.Components.Security) == 0 {
		d.Components.Security = map[string]*SecuritySchemes{ss.TypeName: ss}
		return
	}

	d.Components.Security[ss.TypeName] = ss
}

// AddComponentesSchemasAndExamples responsável por preencher components/schemas e content/contentType/shema, comforme model.
func (d *Doc) AddComponentesSchemasAndExamples(model any, dataType DataTypes) (modelName string) {
	if model == nil {
		return
	}

	_, modelType := d.GetReflectTypeAndValue(model)

	value := d.parseComponentsExamples(model, modelType.Name(), 0, dataType)
	if len(d.Components.Examples) == 0 {
		d.Components.Examples = Examples{}
	}

	d.Components.Examples[modelType.Name()] = &Example{Summary: modelType.Name(), Value: value}
	return d.parseComponentsSchemas(model, modelType.Name(), 0)
}

func (d *Doc) parseComponentsSchemas(model any, ownerModelName string, navigation int) (modelName string) {
	modelValue, modelType := d.GetReflectTypeAndValue(model)
	navigation++

	modelName = modelType.Name()
	if navigation > 3 && modelName == ownerModelName {
		return
	}

	properties := make(map[string]map[string]any, modelValue.NumField())
	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Type().Field(i)

		fieldName := field.Tag.Get("json")

		if fieldName == "" {
			fieldName = field.Name
		}

		fieldType := field.Type

		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}

		switch fieldType.Kind() {
		case reflect.Array, reflect.Slice:
			fieldType = fieldType.Elem()

			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			sub := reflect.New(fieldType).Interface()
			subModelName := d.parseComponentsSchemas(sub, ownerModelName, navigation)

			properties[fieldName] = map[string]any{
				"type": "array",
				"items": map[string]string{
					"$ref": "#/components/schemas/" + subModelName,
				},
			}
			continue

		case reflect.Struct:
			if fieldType != reflect.TypeOf(time.Time{}) {
				sub := reflect.New(fieldType).Interface()
				subModelName := d.parseComponentsSchemas(sub, ownerModelName, navigation)

				properties[fieldName] = map[string]any{
					"$ref": "#/components/schemas/" + subModelName,
				}
				continue
			}
		}

		tp, format, _ := GetDataType(fieldType.Name())
		properties[fieldName] = map[string]any{
			"type":   tp,
			"format": format,
		}
	}

	if d.Components.Schemas == nil {
		d.Components.Schemas = make(map[string]*Schema, 1)
	}

	d.Components.Schemas[modelName] = &Schema{
		Type:       DataTypeObject,
		Properties: properties,
	}
	return
}

func (d *Doc) parseComponentsExamples(model any, ownerModelName string, navigation int, dataType DataTypes) (exampleValue any) {
	modelValue, modelType := d.GetReflectTypeAndValue(model)
	navigation++

	modelName := modelType.Name()
	if navigation > 3 && modelName == ownerModelName {
		return
	}

	exObject := make(map[string]any, modelValue.NumField())

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Type().Field(i)

		fieldName := field.Tag.Get("json")

		if fieldName == "" {
			fieldName = field.Name
		}

		fieldType := field.Type

		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}

		switch fieldType.Kind() {
		case reflect.Array, reflect.Slice:
			fieldType = fieldType.Elem()

			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			sub := reflect.New(fieldType).Interface()
			value := d.parseComponentsExamples(sub, ownerModelName, navigation, DataTypeArray)
			exObject[fieldName] = value
			continue

		case reflect.Struct:
			if fieldType != reflect.TypeOf(time.Time{}) {
				sub := reflect.New(fieldType).Interface()
				value := d.parseComponentsExamples(sub, ownerModelName, navigation, DataTypeObject)
				exObject[fieldName] = value
				continue
			}
		}

		_, _, value := GetDataType(fieldType.Name())

		exObject[fieldName] = value
	}

	exampleValue = exObject

	if dataType == DataTypeArray {
		exampleValue = []map[string]any{exObject}
	}

	return
}

func (d *Doc) GetReflectTypeAndValue(model any) (modelValue reflect.Value, modelType reflect.Type) {
	modelValue = reflect.ValueOf(model)
	modelType = modelValue.Type()

	if modelValue.Kind() == reflect.Pointer {
		modelValue = modelValue.Elem()
		modelType = modelType.Elem()
	}
	return
}
