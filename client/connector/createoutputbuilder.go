package connector

type createOutputBuilder struct {
	createOutput CreateOutput
}

func NewCreateOutputBuilder() *createOutputBuilder {
	return &createOutputBuilder{
		createOutput: CreateOutput{},
	}
}

func (builder *createOutputBuilder) Build() CreateOutput {
	return builder.createOutput
}

func (builder *createOutputBuilder) Uid(uid string) *createOutputBuilder {
	builder.createOutput.Uid = uid
	return builder
}
func (builder *createOutputBuilder) Name(name string) *createOutputBuilder {
	builder.createOutput.Name = name
	return builder
}
func (builder *createOutputBuilder) Status(status string) *createOutputBuilder {
	builder.createOutput.Status = status
	return builder
}
func (builder *createOutputBuilder) State(state string) *createOutputBuilder {
	builder.createOutput.State = state
	return builder
}
func (builder *createOutputBuilder) TenantUid(tenantUid string) *createOutputBuilder {
	builder.createOutput.TenantUid = tenantUid
	return builder
}
func (builder *createOutputBuilder) ServiceConnectivityState(serviceConnectivityState string) *createOutputBuilder {
	builder.createOutput.ServiceConnectivityState = serviceConnectivityState
	return builder
}
func (builder *createOutputBuilder) BootstrapData(bootstrapData string) *createOutputBuilder {
	builder.createOutput.BootstrapData = bootstrapData
	return builder
}
func (builder *createOutputBuilder) CreateRequestOutput(createOutput CreateRequestOutput) *createOutputBuilder {
	builder.createOutput.CreateRequestOutput = &createOutput
	return builder
}
