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
	ParamPath(name string, dataType DataTypes, required bool) PathStructure
	ParamQuery(name string, dataType DataTypes, required bool) PathStructure
	ParamHeader(name string, dataType DataTypes, required bool) PathStructure
	ParamCookie(name string, dataType DataTypes, required bool) PathStructure

	RequestBodyJson(body any, opts ...OptsRequest) PathStructure

	Response(httpStatusCode int, description string) PathStructure
	ResponseBody(contentType string, httpStatusCode int, description string, body any) PathStructure

	ResponseBodyJson(httpStatusCode int, description string, body any) PathStructure
	MethodFunc() (method, pattern string, handlerFn http.HandlerFunc)
	HandleFunc() (methodAndPattern string, handlerFn http.HandlerFunc)
}

func NewDefaultPathStructure(doc *DocJson, method, pattern string, handlerFn http.HandlerFunc, security SecurityType) (p *PathsStructure) {
	tag := "default"

	s := strings.Split(runtime.FuncForPC(reflect.ValueOf(handlerFn).Pointer()).Name(), "/")
	if len(s) > 0 {
		tag = strings.Split(s[len(s)-1], ".")[0]
	}

	p = &PathsStructure{
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

	return
}

// https://swagger.io/docs/specification/paths-and-operations/
type PathsStructure struct {
	Doc         *DocJson         `json:"-"`
	Method      string           `json:"-"`
	Pattern     string           `json:"-"`
	H           http.HandlerFunc `json:"-"`
	Tags        []string         `json:"tags,omitempty"`
	Summ        string           `json:"summary,omitempty"`
	Desc        string           `json:"description,omitempty"`
	Security    []PathSecurity   `json:"security,omitempty"`
	Parameters  []*PathParameter `json:"parameters,omitempty"`
	RequestBody *ResquestBody    `json:"requestBody,omitempty"`
	// A chave representa o http status code (200, 201,..., 400,...)
	Responses map[string]*Response `json:"responses"`
}

type PathSecurity map[string][]string

// https://swagger.io/docs/specification/serialization/
type PathParameter struct {
	Required    bool    `json:"required,omitempty"`
	In          ParamIn `json:"in,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	ParamSchema *Schema `json:"schema,omitempty"`
}

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

func (p *PathsStructure) ParamPath(name string, dataType DataTypes, required bool) PathStructure {
	p.addParameter(ParamPath, name, dataType, required)
	return p
}

func (p *PathsStructure) ParamQuery(name string, dataType DataTypes, required bool) PathStructure {
	p.addParameter(ParamQuery, name, dataType, required)
	return p
}

func (p *PathsStructure) ParamHeader(name string, dataType DataTypes, required bool) PathStructure {
	p.addParameter(ParamHeader, name, dataType, required)
	return p
}

func (p *PathsStructure) ParamCookie(name string, dataType DataTypes, required bool) PathStructure {
	p.addParameter(ParamCookie, name, dataType, required)
	return p
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

	return p.parseBody(content, p.RequestBody.Description, body)
}

func (p *PathsStructure) Response(statusCode int, description string) PathStructure {
	return p.addResponse("", statusCode, description)
}

func (p *PathsStructure) ResponseBody(contentType string, statusCode int, description string, body any) PathStructure {
	return p.addResponse(contentType, statusCode, description, body)
}

func (p *PathsStructure) ResponseBodyJson(statusCode int, description string, body any) PathStructure {
	return p.addResponse("aplication/json", statusCode, description, body)
}

func (p *PathsStructure) addResponse(contentType string, statusCode int, description string, body ...any) PathStructure {
	if strings.TrimSpace(contentType) == "" {
		contentType = "*/*"
	}

	code := fmt.Sprint(statusCode)

	if statusCode < 100 && statusCode > 599 {
		code = "default"
	}

	if len(p.Responses) == 0 {
		p.Responses = make(map[string]*Response, 1)
	}

	resp, ok := p.Responses[code]
	if !ok {
		resp = NewResponse(description)
		p.Responses[code] = resp
	}

	delete(p.Responses, "default")

	if len(body) == 0 {
		return p
	}

	var content *Content
	content, ok = resp.Content[contentType]
	if !ok {
		content = NewContent()
		resp.SetContent(NewContentType(contentType, content))
	}

	return p.parseBody(content, description, body...)
}

func (p *PathsStructure) parseBody(content *Content, description string, body ...any) PathStructure {
	if len(body) == 0 {
		return p
	}

	for _, model := range body {
		modelName, dataType := p.Doc.AddComponentesSchemasAndExamples(model)
		if modelName == "" {
			modelName = strings.ReplaceAll(description, " ", "")
		}
		content.Schemas.AddOneOfRef(modelName, dataType)
		content.AddExamplesRef(modelName)
	}

	return p
}

func (p *PathsStructure) MethodFunc() (method, pattern string, handlerFn http.HandlerFunc) {
	return p.Method, p.Pattern, p.H
}

func (p *PathsStructure) HandleFunc() (methodAndPattern string, handlerFn http.HandlerFunc) {
	return fmt.Sprint(p.Method, p.Pattern), p.H
}

func (p *PathsStructure) addParameter(in ParamIn, name string, sType DataTypes, required bool) {
	p.Parameters = append(p.Parameters, &PathParameter{
		Name:     name,
		In:       in,
		Required: required,
		ParamSchema: &Schema{
			Type: sType,
		},
	})
}
