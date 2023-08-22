package connector

type updateRequestOutputBuilder struct {
	updateRequestOutput UpdateRequestOutput
}

func NewUpdateResponseBuilder() *updateRequestOutputBuilder {
	return &updateRequestOutputBuilder{
		updateRequestOutput: UpdateRequestOutput{},
	}
}

func (builder *updateRequestOutputBuilder) Build() UpdateRequestOutput {
	return builder.updateRequestOutput
}

func (builder *updateRequestOutputBuilder) Uid(uid string) *updateRequestOutputBuilder {
	builder.updateRequestOutput.Uid = uid
	return builder
}
func (builder *updateRequestOutputBuilder) Name(name string) *updateRequestOutputBuilder {
	builder.updateRequestOutput.Name = name
	return builder
}
