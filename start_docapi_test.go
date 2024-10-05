package docapi

import "testing"

func TestNewDocApi(t *testing.T) {
	doc := NewDocApi("localhost:8080/swagger")

	if doc.url != "http://localhost:8080/swagger/" {
		t.Errorf("expected http://localhost:8080/swagger/ but we got %s", doc.url)
	}

	if doc.path != "/swagger/" {
		t.Errorf("expected /swagger/* but we got %s", doc.path)
	}
}
