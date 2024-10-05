// Versão do Swagger 3.0.
//
// https://swagger.io/docs/specification/about/
package docapi

import (
	"strings"
)

type Doc struct {
	Key          string        `json:"-"`
	Version      string        `json:"openapi"`
	Servers      []Servers     `json:"servers,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
	Info         *Info         `json:"info,omitempty"`
	// a chave é o path do endpoint
	Paths      map[string]Path `json:"paths,omitempty"`
	Components *Components     `json:"components,omitempty"`
}

type Servers struct {
	URL string `json:"url"`
}

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Info struct {
	Description string   `json:"description,omitempty"`
	Version     string   `json:"version,omitempty"`
	Title       string   `json:"title,omitempty"`
	Contact     *Contact `json:"contact,omitempty"`
	License     *License `json:"license,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

func (j *Doc) AddServer(url ...string) {
	for _, v := range url {
		if !j.serverIsPresent(v) {
			j.Servers = append(j.Servers, Servers{URL: v})
		}
	}
}

func (j *Doc) serverIsPresent(url string) (ok bool) {
	for _, v := range j.Servers {
		if ok = v.URL == url; ok {
			return
		}
	}
	return
}

func (j *Doc) AddPath(method, pattern string, pathStructure *PathsStructure) {
	method = strings.ToLower(method)
	// A raiz do path é a url e dentro contém os métodos get, post, put...
	// Se localizar a url, então adiciona o método.
	if len(j.Paths) == 0 {
		j.Paths = map[string]Path{pattern: {method: pathStructure}}
		return
	}

	if paths, ok := j.Paths[pattern]; ok {
		paths[method] = pathStructure
		j.Paths[pattern] = paths

	} else {
		j.Paths[pattern] = Path{method: pathStructure}
	}
}
