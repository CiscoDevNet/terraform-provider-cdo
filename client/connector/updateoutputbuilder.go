package connector

type updateOutputBuilder struct {
	updateOutput UpdateOutput
}

func NewUpdateOutputBuilder() *updateOutputBuilder {
	return &updateOutputBuilder{
		updateOutput: UpdateOutput{},
	}
}

func (builder *updateOutputBuilder) Build() UpdateOutput {
	return builder.updateOutput
}

func (builder *updateOutputBuilder) Uid(uid string) *updateOutputBuilder {
	builder.updateOutput.Uid = uid
	return builder
}
func (builder *updateOutputBuilder) Name(name string) *updateOutputBuilder {
	builder.updateOutput.Name = name
	return builder
}

func (builder *updateOutputBuilder) BootstrapData(bootstrapData string) *updateOutputBuilder {
	builder.updateOutput.BootstrapData = bootstrapData
	return builder
}
func (builder *updateOutputBuilder) UpdateRequestOutput(updateOutput UpdateRequestOutput) *updateOutputBuilder {
	builder.updateOutput.UpdateRequestOutput = &updateOutput
	return builder
}
