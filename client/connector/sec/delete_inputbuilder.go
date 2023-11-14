package sec

type DeleteInputBuilder struct {
	deleteInput *DeleteInput
}

func NewDeleteInputBuilder() *DeleteInputBuilder {
	deleteInput := &DeleteInput{}
	b := &DeleteInputBuilder{deleteInput: deleteInput}
	return b
}

func (b *DeleteInputBuilder) Uid(uid string) *DeleteInputBuilder {
	b.deleteInput.Uid = uid
	return b
}

func (b *DeleteInputBuilder) Build() DeleteInput {
	return *b.deleteInput
}
