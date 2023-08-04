package sdc

type createRequestOutputBuilder struct {
	createRequestOutput CreateRequestOutput
}

func NewCreateResponseBuilder() *createRequestOutputBuilder {
	return &createRequestOutputBuilder{
		createRequestOutput: CreateRequestOutput{},
	}
}

func (builder *createRequestOutputBuilder) Build() CreateRequestOutput {
	return builder.createRequestOutput
}

func (builder *createRequestOutputBuilder) Uid(uid string) *createRequestOutputBuilder {
	builder.createRequestOutput.Uid = uid
	return builder
}
func (builder *createRequestOutputBuilder) Name(name string) *createRequestOutputBuilder {
	builder.createRequestOutput.Name = name
	return builder
}
func (builder *createRequestOutputBuilder) Status(status string) *createRequestOutputBuilder {
	builder.createRequestOutput.Status = status
	return builder
}
func (builder *createRequestOutputBuilder) State(state string) *createRequestOutputBuilder {
	builder.createRequestOutput.State = state
	return builder
}
func (builder *createRequestOutputBuilder) TenantUid(tenantUid string) *createRequestOutputBuilder {
	builder.createRequestOutput.TenantUid = tenantUid
	return builder
}
func (builder *createRequestOutputBuilder) ServiceConnectivityState(serviceConnectivityState string) *createRequestOutputBuilder {
	builder.createRequestOutput.ServiceConnectivityState = serviceConnectivityState
	return builder
}
