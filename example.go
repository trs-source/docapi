package docapi

type Examples map[string]*Example

type Example struct {
	// Tokens usado para manter ordenado o Value, conforme sequência de campos da struct modelo (dto).
	// No GET do endpoint doc.json é removido o token do payload.
	Tokens  [][]byte `json:"-"`
	Summary string   `json:"summary,omitempty"`
	Value   any      `json:"value,omitempty"`
	Ref     string   `json:"$ref,omitempty"`
	// TypeName usado quando o modelo(dto) foi criado com reflect, neste caso não tem o nome da struct.
	TypeName string `json:"-"`
}

type OptsExample func(*Example)

func WithExampleSummary(summary string) OptsExample {
	return func(e *Example) {
		e.Summary = summary
	}
}

// WithTypeName é usado quando a struct é criada via reflect, neste caso não tem o nome dela.
func WithTypeName(typeName string) OptsExample {
	return func(e *Example) {
		e.TypeName = typeName
	}
}
