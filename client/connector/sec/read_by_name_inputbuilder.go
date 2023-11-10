package sec

type ReadByNameInputBuilder struct {
	readByNameInput *ReadByNameInput
}

func NewReadByNameInputBuilder() *ReadByNameInputBuilder {
	readByNameInput := &ReadByNameInput{}
	b := &ReadByNameInputBuilder{readByNameInput: readByNameInput}
	return b
}

func (b *ReadByNameInputBuilder) Name(name string) *ReadByNameInputBuilder {
	b.readByNameInput.Name = name
	return b
}

func (b *ReadByNameInputBuilder) Build() ReadByNameInput {
	return *b.readByNameInput
}
