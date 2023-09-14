package device

type ReadSpecificInputBuilder struct {
	readSpecificInput *ReadSpecificInput
}

func NewReadSpecificInputBuilder() *ReadSpecificInputBuilder {
	readSpecificInput := &ReadSpecificInput{}
	b := &ReadSpecificInputBuilder{readSpecificInput: readSpecificInput}
	return b
}

func (b *ReadSpecificInputBuilder) Uid(uid string) *ReadSpecificInputBuilder {
	b.readSpecificInput.Uid = uid
	return b
}

func (b *ReadSpecificInputBuilder) Build() ReadSpecificInput {
	return *b.readSpecificInput
}
