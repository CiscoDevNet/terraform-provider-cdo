package sec

type ReadInputBuilder struct {
	readInput *ReadInput
}

func NewReadInputBuilder() *ReadInputBuilder {
	readInput := &ReadInput{}
	b := &ReadInputBuilder{readInput: readInput}
	return b
}

func (b *ReadInputBuilder) Uid(uid string) *ReadInputBuilder {
	b.readInput.Uid = uid
	return b
}

func (b *ReadInputBuilder) Build() ReadInput {
	return *b.readInput
}
