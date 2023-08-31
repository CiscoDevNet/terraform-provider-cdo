package fmcappliance

type UpdateInputBuilder struct {
	updateInput *UpdateInput
}

func NewUpdateInputBuilder() *UpdateInputBuilder {
	updateInput := &UpdateInput{}
	b := &UpdateInputBuilder{updateInput: updateInput}
	return b
}

func (b *UpdateInputBuilder) FmcApplianceUid(fmcApplianceUid string) *UpdateInputBuilder {
	b.updateInput.FmcApplianceUid = fmcApplianceUid
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
