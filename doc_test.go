package docapi

import "testing"

type Model struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Childer []Model `json:"childer"`
}

func TestParseComponentsExamples(t *testing.T) {
	doc := Session().NewDoc("/test")

	doc.AddComponentesSchemasAndExamples(Model{})

}
