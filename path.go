package docapi

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type PathStructure interface {
	Tag(string) PathStructure
	Summary(string) PathStructure
	Description(string) PathStructure
	ParamPath(name string, dataType DataType, opts ...OptsParameter) PathStructure
	ParamQuery(name string, dataType DataType, opts ...OptsParameter) PathStructure
	ParamHeader(name string, dataType DataType, opts ...OptsParameter) PathStructure
	ParamCookie(name string, dataType DataType, opts ...OptsParameter) PathStructure
	RequestBodyJson(body any, opts ...OptsRequest) PathStructure
	Response(httpStatusCode int, description string) PathStructure
	ResponseBody(contentType string, httpStatusCode int, description string, body any, opts ...OptsExample) PathStructure
	ResponseBodyJson(httpStatusCode int, description string, body any, opts ...OptsExample) PathStructure
	MethodFunc() (method, pattern string, handlerFn http.HandlerFunc)
	// HandleFunc retornar o método e path na mesma string (ex.: GET /busca-os), com intuito de ser usado no net/http (nativo).
	HandleFunc() (methodAndPattern string, handlerFn http.HandlerFunc)
}

func NewDefaultPathStructure(doc *Doc, method, pattern string, handlerFn http.HandlerFunc, security SecurityType) PathStructure {
	tag := "default"

	// Obtém o nome do arquivo que está o controller
	s := strings.Split(runtime.FuncForPC(reflect.ValueOf(handlerFn).Pointer()).Name(), "/")
	if len(s) > 0 {
		tag = strings.Split(s[len(s)-1], ".")[0]
	}

	p := &PathsStructure{
		Doc:       doc,
		Method:    method,
		Pattern:   pattern,
		H:         handlerFn,
		Tags:      []string{tag},
		Responses: map[string]*Response{"default": {Description: "Default"}},
	}

	if security != SecurityNone {
		p.Security = append(p.Security, PathSecurity{security.String(): []string{}})
	}

	doc.AddPath(method, pattern, p)
	return p
}

// A chave é o método http (delete, get, patch, post...)
type Path map[string]*PathsStructure

var _ PathStructure = (*PathsStructure)(nil)

// https://swagger.io/docs/specification/paths-and-operations/
type PathsStructure struct {
	Doc         *Doc             `json:"-"`
	Method      string           `json:"-"`
	Pattern     string           `json:"-"`
	H           http.HandlerFunc `json:"-"`
	Tags        []string         `json:"tags,omitempty"`
	Summ        string           `json:"summary,omitempty"`
	Desc        string           `json:"description,omitempty"`
	Security    []PathSecurity   `json:"security,omitempty"`
	Parameters  []*Parameter     `json:"parameters,omitempty"`
	RequestBody *ResquestBody    `json:"requestBody,omitempty"`
	// A chave representa o http status code (200, 201,..., 400,...)
	Responses map[string]*Response `json:"responses"`
}

type PathSecurity map[string][]string

func (p *PathsStructure) Tag(tag string) PathStructure {
	p.Tags = []string{tag}
	return p
}

func (p *PathsStructure) Summary(summary string) PathStructure {
	p.Summ = summary
	return p
}

func (p *PathsStructure) Description(description string) PathStructure {
	p.Desc = description
	return p
}

func (p *PathsStructure) ParamPath(name string, dataType DataType, opts ...OptsParameter) PathStructure {
	p.addParameter(ParamPath, name, dataType, opts...)
	return p
}

func (p *PathsStructure) ParamQuery(name string, dataType DataType, opts ...OptsParameter) PathStructure {
	p.addParameter(ParamQuery, name, dataType, opts...)
	return p
}

func (p *PathsStructure) ParamHeader(name string, dataType DataType, opts ...OptsParameter) PathStructure {
	p.addParameter(ParamHeader, name, dataType, opts...)
	return p
}

func (p *PathsStructure) ParamCookie(name string, dataType DataType, opts ...OptsParameter) PathStructure {
	p.addParameter(ParamCookie, name, dataType, opts...)
	return p
}

func (p *PathsStructure) addParameter(in ParamIn, name string, sType DataType, opts ...OptsParameter) {
	param := &Parameter{
		Name: name,
		In:   in,
		ParamSchema: &Schema{
			Type: sType,
		},
	}

	for _, fn := range opts {
		fn(param)
	}

	p.Parameters = append(p.Parameters, param)
}

func (p *PathsStructure) RequestBodyJson(body any, opts ...OptsRequest) PathStructure {
	return p.setRequest("aplication/json", body, opts...)
}

func (p *PathsStructure) setRequest(contentType string, body any, opts ...OptsRequest) PathStructure {
	if strings.TrimSpace(contentType) == "" {
		contentType = "*/*"
	}

	if p.RequestBody == nil {
		p.RequestBody = &ResquestBody{
			Content: NewContentType(contentType, NewContent()),
		}
	}

	content, ok := p.RequestBody.Content[contentType]
	if !ok {
		content = NewContent()
		p.RequestBody.Content = NewContentType(contentType, content)
	}

	for _, fn := range opts {
		fn(p.RequestBody)
	}

	if body == nil {
		content.Schemas.Type = DataTypeString
		return p
	}

	var optsExample []OptsExample
	if p.RequestBody.exempleSummary != "" {
		optsExample = append(optsExample, WithExampleSummary(p.RequestBody.exempleSummary))
	}

	if p.RequestBody.typeName != "" {
		optsExample = append(optsExample, WithTypeName(p.RequestBody.typeName))
	}

	return p.parseBody(content, body, optsExample...)
}

func (p *PathsStructure) Response(statusCode int, description string) PathStructure {
	return p.addResponse("", statusCode, description, nil)
}

func (p *PathsStructure) ResponseBody(contentType string, statusCode int, description string, body any, opts ...OptsExample) PathStructure {
	return p.addResponse(contentType, statusCode, description, body, opts...)
}

func (p *PathsStructure) ResponseBodyJson(statusCode int, description string, body any, opts ...OptsExample) PathStructure {
	return p.addResponse("aplication/json", statusCode, description, body, opts...)
}

func (p *PathsStructure) addResponse(contentType string, statusCode int, description string, body any, opts ...OptsExample) PathStructure {
	if strings.TrimSpace(contentType) == "" {
		contentType = "*/*"
	}

	code := fmt.Sprint(statusCode)

	if len(p.Responses) == 0 {
		p.Responses = make(map[string]*Response, 1)
	}

	// deleta o default que foi criado na estrutura inicial.
	delete(p.Responses, "default")

	if statusCode < 100 && statusCode > 599 {
		code = "default"
	}

	resp, ok := p.Responses[code]
	if !ok {
		resp = NewResponse(description)
		p.Responses[code] = resp
	}

	if body == nil {
		return p
	}

	var content *Content
	content, ok = resp.Content[contentType]
	if !ok {
		content = NewContent()
		resp.SetContent(NewContentType(contentType, content))
	}

	return p.parseBody(content, body, opts...)
}

func (p *PathsStructure) parseBody(content *Content, body any, opts ...OptsExample) PathStructure {
	if body == nil {
		return p
	}

	modelValue := reflect.ValueOf(body)
	modelType := modelValue.Type()

	if modelValue.Kind() == reflect.Pointer {
		modelValue = modelValue.Elem()
		modelType = modelType.Elem()
	}

	var dataType DataType

	switch modelType.Kind() {
	case reflect.Slice, reflect.Array:
		dataType = DataTypeArray
		modelType = modelType.Elem()
		if modelType.Kind() == reflect.Pointer {
			modelType = modelType.Elem()
		}
		modelValue = reflect.ValueOf(reflect.New(modelType).Interface()).Elem()

	case reflect.Struct:
		dataType = DataTypeObject
	default:
		return p
	}

	modelName := p.Doc.Components.AddSchemasAndExamples(modelValue, modelType, dataType, opts...)

	content.Schemas.AddOneOfRef(modelName, dataType)
	content.AddExamplesRef(modelName)

	return p
}

func (p *PathsStructure) MethodFunc() (method, pattern string, handlerFn http.HandlerFunc) {
	return p.Method, p.Pattern, p.H
}

func (p *PathsStructure) HandleFunc() (methodAndPattern string, handlerFn http.HandlerFunc) {
	return fmt.Sprint(strings.ToUpper(p.Method), " ", p.Pattern), p.H
}
