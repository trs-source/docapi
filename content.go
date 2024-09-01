package docapi

type ContentType map[string]*Content

type Content struct {
	Schemas  *Schema  `json:"schema,omitempty"`
	Examples Examples `json:"examples,omitempty"`
}

func (c *Content) AddExamplesRef(modelName string) {
	ref := "#/components/examples/" + modelName
	if len(c.Examples) == 0 {
		c.Examples = Examples{
			modelName: &Example{
				Ref: ref},
		}
		return
	}

	c.Examples[modelName] = &Example{Ref: ref}
}

func NewContent() *Content {
	return &Content{Schemas: &Schema{}}
}

func NewContentType(contentType string, content *Content) ContentType {
	return ContentType{contentType: content}
}
