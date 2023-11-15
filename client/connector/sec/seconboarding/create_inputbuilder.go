package seconboarding

type CreateInputBuilder struct {
	createInput *CreateInput
}

func NewCreateInputBuilder() *CreateInputBuilder {
	createInput := &CreateInput{}
	b := &CreateInputBuilder{createInput: createInput}
	return b
}

func (b *CreateInputBuilder) Name(name string) *CreateInputBuilder {
	b.createInput.Name = name
	return b
}

func (b *CreateInputBuilder) Build() CreateInput {
	return *b.createInput
}
