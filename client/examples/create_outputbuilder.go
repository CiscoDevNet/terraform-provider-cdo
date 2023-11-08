package examples

type CreateInputBuilder struct {
	createInput *CreateInput
}

func NewCreateInputBuilder() *CreateInputBuilder {
	createInput := &CreateInput{}
	b := &CreateInputBuilder{createInput: createInput}
	return b
}

func (b *CreateInputBuilder) Uid(uid string) *CreateInputBuilder {
	b.createInput.Uid = uid
	return b
}

func (b *CreateInputBuilder) Build() CreateInput {
	return *b.createInput
}
