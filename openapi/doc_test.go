package openapi

import "testing"

type Model struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Childer []Model `json:"childer"`
}

func TestParseComponentsExamples(t *testing.T) {
	doc := Mapping().NewDoc("/test")
	_, m := doc.GetReflectTypeAndValue(Model{})

	value := doc.parseComponentsExamples(Model{}, m.Name(), 0, DataTypeObject)

	if value, ok := value.(map[string]any); !ok {
		t.Error("example tag failed")
	} else {
		_ = value
	}

	value = doc.parseComponentsExamples(Model{}, m.Name(), 0, DataTypeArray)

	if value, ok := value.([]map[string]any); !ok {
		t.Error("example tag failed")
	} else {
		_ = value
	}
}
