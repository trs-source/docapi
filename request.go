package docapi

// https://swagger.io/docs/specification/describing-request-body/
type ResquestBody struct {
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Content     ContentType `json:"content,omitempty"`
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

func WithRequired(required bool) OptsRequest {
	return func(r *ResquestBody) {
		r.Required = true
	}
}
