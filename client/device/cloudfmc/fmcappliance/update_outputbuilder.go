package fmcappliance

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

func (b *updateRequestBodyBuilder) StateMachineContext(stateMachineContext *map[string]string) *updateRequestBodyBuilder {
	b.updateRequestBody.StateMachineContext = stateMachineContext
	return b
}

func (b *updateRequestBodyBuilder) Uid(uid string) *updateRequestBodyBuilder {
	b.updateRequestBody.Uid = uid
	return b
}

func (b *updateRequestBodyBuilder) Build() updateRequestBody {
	return *b.updateRequestBody
}
