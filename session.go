package docapi

import (
	"log/slog"
	"sync"
)

var (
	once   sync.Once
	mapper *Mapper
)

func Session() *Mapper {
	once.Do(
		func() {
			mapper = &Mapper{Docs{}}
		})
	return mapper
}

type Mapper struct {
	Docs Docs
}

// NewDoc responsável por criar a configuração padrão para o doc.json.
//
// key: Representa o path que irá acessar o doc.json. Ex.: /swagger/doc.json
func (m *Mapper) NewDoc(key string) *DocJson {
	doc := &DocJson{
		Key:     key,
		Version: "3.0.1",
		Info: &Info{
			Title:   "DocApi",
			Version: "1.0",
		},
		Components: &Components{},
	}
	m.Docs[key] = doc
	return doc
}

// FindDocByPathDocJson localiza o doc que está vinculado a key.
//
// path: Representa o path que contém o final doc.json.
// Ex.: URL = http://localhost:8080/swagger: path = /swagger/doc.json.
func (m *Mapper) FindDocByPathDocJson(path string) (doc *DocJson, ok bool) {
	doc, ok = m.Docs[path]
	if ok {
		return
	}

	slog.Error("[DocApi] doc.json not found by key.", "key", path)
	return
}
