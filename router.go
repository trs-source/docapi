package docapi

import (
	"net/http"
)

type Router struct {
	key      string
	security SecurityType
}

func newRouter(key string, security SecurityType) *Router {
	return &Router{key: key, security: security}
}

func (o *Router) Connect(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "connect", pattern, handlerFn, o.security)
}

func (o *Router) Delete(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "delete", pattern, handlerFn, o.security)
}

func (o *Router) Get(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "get", pattern, handlerFn, o.security)
}

func (o *Router) Head(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "head", pattern, handlerFn, o.security)
}

func (o *Router) Patch(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "patch", pattern, handlerFn, o.security)
}

func (o *Router) Post(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "post", pattern, handlerFn, o.security)
}

func (o *Router) Put(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "put", pattern, handlerFn, o.security)
}

func (o *Router) Options(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "options", pattern, handlerFn, o.security)
}

func (o *Router) Trace(pattern string, handlerFn http.HandlerFunc) PathStructure {
	return NewDefaultPathStructure(o.key, "trace", pattern, handlerFn, o.security)
}
