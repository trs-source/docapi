package docapi

// https://swagger.io/docs/specification/describing-request-body/
type ResquestBody struct {
	Description    string      `json:"description,omitempty"`
	Required       bool        `json:"required,omitempty"`
	Content        ContentType `json:"content,omitempty"`
	exempleSummary string      `json:"-"`
	typeName       string      `json:"-"`
}

func NewRequest(description string) *ResquestBody {
	return &ResquestBody{
		Description: description,
	}
}

func (r *ResquestBody) SetContent(contentType ContentType) {
	r.Content = contentType
}

type OptsRequest func(*ResquestBody)

func WithDescription(description string) OptsRequest {
	return func(r *ResquestBody) {
		r.Description = description
	}
}

func WithRequired() OptsRequest {
	return func(r *ResquestBody) {
		r.Required = true
	}
}

func WithReqExampleSummary(summary string) OptsRequest {
	return func(r *ResquestBody) {
		r.exempleSummary = summary
	}
}

// WithReqTypeName é usado quando a struct é criada via reflect, neste caso não tem o nome dela.
func WithReqTypeName(typeName string) OptsRequest {
	return func(r *ResquestBody) {
		r.typeName = typeName
	}
}
