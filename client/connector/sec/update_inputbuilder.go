package sec

type UpdateInputBuilder struct {
	updateInput *UpdateInput
}

func NewUpdateInputBuilder() *UpdateInputBuilder {
	updateInput := &UpdateInput{}
	b := &UpdateInputBuilder{updateInput: updateInput}
	return b
}

func (b *UpdateInputBuilder) Build() UpdateInput {
	return *b.updateInput
}
