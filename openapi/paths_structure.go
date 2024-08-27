package openapi

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
	ParamURL(name string, dataType DataTypes, required bool) PathStructure
	ParamQuery(name string, dataType DataTypes, required bool) PathStructure
	Response(httpStatusCode int, description string) PathStructure
	ResponseBody(contentType string, httpStatusCode int, description string, dataType DataTypes, body ...any) PathStructure
	ResponseTextPlain(httpStatusCode int, description string) PathStructure
	ResponseObjectBodyJson(httpStatusCode int, description string, body ...any) PathStructure
	ResponseArrayBodyJson(httpStatusCode int, description string, body ...any) PathStructure

	HandlerFn() (method, pattern string, handlerFn http.HandlerFunc)
}

func NewDefaultPathStructure(key string, method, pattern string, handlerFn http.HandlerFunc, security SecurityType) (p *PathsStructure) {
	tag := "default"

	s := strings.Split(runtime.FuncForPC(reflect.ValueOf(handlerFn).Pointer()).Name(), "/")
	if len(s) > 0 {
		tag = strings.Split(s[len(s)-1], ".")[0]
	}

	p = &PathsStructure{
		Key:       key,
		Method:    method,
		Pattern:   pattern,
		H:         handlerFn,
		Tags:      []string{tag},
		Responses: map[string]*Response{"200": {Description: "OK"}},
	}

	if security != SecurityNone {
		p.Security = append(p.Security, PathSecurity{security.String(): []string{}})
	}

	if doc, ok := Mapping().FindDocByKey(key); ok {
		doc.AddPath(method, pattern, p)
	}

	return
}

// https://swagger.io/docs/specification/paths-and-operations/
type PathsStructure struct {
	Key         string           `json:"-"`
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

// https://swagger.io/docs/specification/describing-request-body/
type ResquestBody struct {
	Required bool         `json:"required,omitempty"`
	Content  *ContentType `json:"content,omitempty"`
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

func (p *PathsStructure) ParamURL(name string, dataType DataTypes, required bool) PathStructure {
	p.addParameter(ParamPath, name, dataType, required)
	return p
}

func (p *PathsStructure) ParamQuery(name string, dataType DataTypes, required bool) PathStructure {
	p.addParameter(ParamQuery, name, dataType, required)
	return p
}

func (p *PathsStructure) Response(statusCode int, description string) PathStructure {
	return p.addResponse("", statusCode, description, DataTypeString)
}

func (p *PathsStructure) ResponseBody(contentType string, statusCode int, description string, dataType DataTypes, body ...any) PathStructure {
	return p.addResponse(contentType, statusCode, description, dataType, body...)
}

func (p *PathsStructure) ResponseTextPlain(statusCode int, description string) PathStructure {
	return p.addResponse("text/plain", statusCode, description, DataTypeString)
}

func (p *PathsStructure) ResponseObjectBodyJson(statusCode int, description string, body ...any) PathStructure {
	return p.addResponse("aplication/json", statusCode, description, DataTypeObject, body...)
}

func (p *PathsStructure) ResponseArrayBodyJson(statusCode int, description string, body ...any) PathStructure {
	return p.addResponse("aplication/json", statusCode, description, DataTypeArray, body...)
}

func (p *PathsStructure) addResponse(headerContentType string, statusCode int, description string, dataType DataTypes, body ...any) PathStructure {
	if strings.TrimSpace(headerContentType) == "" {
		headerContentType = "*/*"
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

	if len(body) == 0 {
		return p
	}

	doc, ok := Mapping().FindDocByKey(p.Key)
	if !ok {
		return p
	}

	for _, model := range body {
		modelName := doc.AddComponentesSchemasAndExamples(model, dataType)

		if !ok {
			resp.SetContent(NewContentType(dataType, headerContentType, modelName))
			ok = true
			continue
		}

		if _, okc := resp.Content[headerContentType]; okc {
			resp.Content[headerContentType].Schemas.AddOneOfRef(modelName, dataType)
			resp.Content[headerContentType].AddExamplesRef(modelName)
		} else {
			resp.SetContent(NewContentType(dataType, headerContentType, modelName))
		}

	}
	return p
}

func (p *PathsStructure) HandlerFn() (method, pattern string, handlerFn http.HandlerFunc) {
	return p.Method, p.Pattern, p.H
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
