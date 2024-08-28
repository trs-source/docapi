package docapi

import "testing"

func TestNewDocApi(t *testing.T) {
	doc := NewDocApi("localhost:8080/swagger")

	if doc.url != "http://localhost:8080/swagger/" {
		t.Errorf("expected http://localhost:8080/swagger/ but we got %s", doc.url)
	}

	if doc.pattern != "/swagger/*" {
		t.Errorf("expected /swagger/* but we got %s", doc.pattern)
	}

	if doc.key != "/swagger/doc.json" {
		t.Errorf("expected /swagger/doc.json but we got %s", doc.key)
	}
}
