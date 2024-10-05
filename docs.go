package docapi

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
)

var (
	once sync.Once
	doc  *Docs
)

func GetDocs() *Docs {
	once.Do(
		func() {
			doc = &Docs{make(map[string]*Doc)}
		})
	return doc
}

type Docs struct {
	Docs map[string]*Doc
}

// NewDoc responsável por criar a configuração padrão para o doc.json.
//
// key: Representa o path que irá acessar o doc.json. Ex.: /swagger/doc.json
func (d *Docs) NewDoc(pathDocJson string) *Doc {
	doc := &Doc{
		Key:     pathDocJson,
		Version: "3.0.1",
		Info: &Info{
			Title:   "DocApi",
			Version: "1.0",
		},
		Components: &Components{},
	}
	d.Docs[pathDocJson] = doc
	return doc
}

// FindDocJSONByPath localiza o docJSON que está vinculado ao path.
//
// path: Ex.: URL = http://localhost:8080/swagger: path = /swagger/
func (d *Docs) FindDocJSONByPath(path string) (doc *Doc, ok bool) {
	path = strings.TrimSuffix(path, "/")
	doc, ok = d.Docs[path+"/doc.json"]
	if ok {
		return
	}

	slog.Error("[DocApi] doc.json not found by path.", "path", path)
	return
}

func (d *Docs) GetJSON(rURLPath string) (response []byte) {
	doc, ok := d.FindDocJSONByPath(strings.TrimSuffix(rURLPath, "/doc.json"))
	if !ok {
		return
	}

	var err error
	response, err = json.Marshal(doc)
	if err != nil {
		slog.Error("error when creating Swagger doc.json file.", "error", err.Error())
		return
	}

	if doc.Components == nil {
		return
	}

	for _, v := range doc.Components.Examples {
		for _, token := range v.Tokens {
			response = bytes.ReplaceAll(response, token, []byte(""))
		}

	}
	return
}
