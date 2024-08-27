package openapi

import (
	"log/slog"
	"sync"
)

var (
	once   sync.Once
	mapper *Mapper
)

func Mapping() *Mapper {
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
func (m *Mapper) NewDoc(key string) *Doc {
	doc := &Doc{
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

func (m *Mapper) FindDocByKey(key string) (doc *Doc, ok bool) {
	doc, ok = m.Docs[key]
	if ok {
		return
	}

	slog.Error("[DocApi] doc.json not found by url.", "Url", key)
	return
}
