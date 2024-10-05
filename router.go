package docapi

import (
	"net/http"
)

type Router struct {
	security SecurityType
	document *Doc
}

func newRouter(doc *Doc, security SecurityType) Router {
	return Router{document: doc, security: security}
}

func (o Router) Connect(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "connect", pattern, handlerFn, o.security)
}

func (o Router) Delete(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "delete", pattern, handlerFn, o.security)
}

func (o Router) Get(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "get", pattern, handlerFn, o.security)
}

func (o Router) Head(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "head", pattern, handlerFn, o.security)
}

func (o Router) Patch(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "patch", pattern, handlerFn, o.security)
}

func (o Router) Post(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "post", pattern, handlerFn, o.security)
}

func (o Router) Put(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "put", pattern, handlerFn, o.security)
}

func (o Router) Options(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "options", pattern, handlerFn, o.security)
}

func (o Router) Trace(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.document, "trace", pattern, handlerFn, o.security)
}
