package docapi

import (
	"net/http"
)

type Router struct {
	key      string
	security SecurityType
	doc      *DocJson
}

func newRouter(doc *DocJson, security SecurityType) *Router {
	return &Router{doc: doc, security: security}
}

func (o *Router) Connect(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "connect", pattern, handlerFn, o.security)
}

func (o *Router) Delete(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "delete", pattern, handlerFn, o.security)
}

func (o *Router) Get(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "get", pattern, handlerFn, o.security)
}

func (o *Router) Head(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "head", pattern, handlerFn, o.security)
}

func (o *Router) Patch(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "patch", pattern, handlerFn, o.security)
}

func (o *Router) Post(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "post", pattern, handlerFn, o.security)
}

func (o *Router) Put(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "put", pattern, handlerFn, o.security)
}

func (o *Router) Options(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "options", pattern, handlerFn, o.security)
}

func (o *Router) Trace(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.doc, "trace", pattern, handlerFn, o.security)
}
