package openapi

// https://swagger.io/docs/specification/describing-responses/
type Response struct {
	Description string      `json:"description"`
	Content     ContentType `json:"content,omitempty"`
}

func NewResponse(description string) *Response {
	return &Response{
		Description: description,
	}
}

func (r *Response) SetContent(contentType ContentType) {
	r.Content = contentType
}
