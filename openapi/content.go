package openapi

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

func NewContentType(dataType DataTypes, contentType, modelName string) ContentType {
	schema := NewSchemaForResponsesContent(modelName, dataType)
	content := &Content{Schemas: schema}
	content.AddExamplesRef(modelName)
	return ContentType{contentType: content}
}
