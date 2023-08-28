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

type UpdateOutputBuilder struct {
	updateOutput *UpdateOutput
}

func NewUpdateOutputBuilder() *UpdateOutputBuilder {
	updateOutput := &UpdateOutput{}
	b := &UpdateOutputBuilder{updateOutput: updateOutput}
	return b
}

func (b *UpdateOutputBuilder) Uid(uid string) *UpdateOutputBuilder {
	b.updateOutput.Uid = uid
	return b
}

func (b *UpdateOutputBuilder) State(state string) *UpdateOutputBuilder {
	b.updateOutput.State = state
	return b
}

func (b *UpdateOutputBuilder) DomainUid(domainUid string) *UpdateOutputBuilder {
	b.updateOutput.DomainUid = domainUid
	return b
}

func (b *UpdateOutputBuilder) Build() UpdateOutput {
	return *b.updateOutput
}

type updateRequestBodyBuilder struct {
	updateRequestBody *updateRequestBody
}

func newUpdateRequestBodyBuilder() *updateRequestBodyBuilder {
	updateRequestBody := &updateRequestBody{}
	b := &updateRequestBodyBuilder{updateRequestBody: updateRequestBody}
	return b
}

func (b *updateRequestBodyBuilder) QueueTriggerState(queueTriggerState string) *updateRequestBodyBuilder {
	b.updateRequestBody.QueueTriggerState = queueTriggerState
	return b
}

func (b *updateRequestBodyBuilder) StateMachineContext(stateMachineContext map[string]string) *updateRequestBodyBuilder {
	b.updateRequestBody.StateMachineContext = stateMachineContext
	return b
}

func (b *updateRequestBodyBuilder) Build() updateRequestBody {
	return *b.updateRequestBody
}
