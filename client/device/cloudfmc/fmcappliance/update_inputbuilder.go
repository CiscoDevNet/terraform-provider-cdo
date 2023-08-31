package fmcappliance

type UpdateInputBuilder struct {
	updateInput *UpdateInput
}

func NewUpdateInputBuilder() *UpdateInputBuilder {
	updateInput := &UpdateInput{}
	b := &UpdateInputBuilder{updateInput: updateInput}
	return b
}

func (b *UpdateInputBuilder) FmcSpecificUid(fmcSpecificUid string) *UpdateInputBuilder {
	b.updateInput.FmcSpecificUid = fmcSpecificUid
	return b
}

func (b *UpdateInputBuilder) QueueTriggerState(queueTriggerState string) *UpdateInputBuilder {
	b.updateInput.QueueTriggerState = queueTriggerState
	return b
}

func (b *UpdateInputBuilder) StateMachineContext(stateMachineContext map[string]string) *UpdateInputBuilder {
	b.updateInput.StateMachineContext = stateMachineContext
	return b
}

func (b *UpdateInputBuilder) Build() UpdateInput {
	return *b.updateInput
}
