package openapi

type Examples map[string]*Example

type Example struct {
	Summary string `json:"summary,omitempty"`
	Value   any    `json:"value,omitempty"`
	Ref     string `json:"$ref,omitempty"`
}
