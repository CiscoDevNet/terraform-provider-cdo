package user

type getTokenOutputBuilder struct {
	getTokenOutputBuilder GetTokenOutput
}

func NewGetTokenOutputBuilder() *getTokenOutputBuilder {
	return &getTokenOutputBuilder{
		getTokenOutputBuilder: GetTokenOutput{},
	}
}

func (builder *getTokenOutputBuilder) Build() GetTokenOutput {
	return builder.getTokenOutputBuilder
}

func (builder *getTokenOutputBuilder) TenantUid(tenantUid string) *getTokenOutputBuilder {
	builder.getTokenOutputBuilder.TenantUid = tenantUid
	return builder
}

func (builder *getTokenOutputBuilder) TenantName(tenantName string) *getTokenOutputBuilder {
	builder.getTokenOutputBuilder.TenantName = tenantName
	return builder
}

func (builder *getTokenOutputBuilder) AccessToken(accessToken string) *getTokenOutputBuilder {
	builder.getTokenOutputBuilder.AccessToken = accessToken
	return builder
}

func (builder *getTokenOutputBuilder) RefreshToken(refreshToken string) *getTokenOutputBuilder {
	builder.getTokenOutputBuilder.RefreshToken = refreshToken
	return builder
}

func (builder *getTokenOutputBuilder) TokenType(tokenType string) *getTokenOutputBuilder {
	builder.getTokenOutputBuilder.TokenType = tokenType
	return builder
}

func (builder *getTokenOutputBuilder) Scope(scope string) *getTokenOutputBuilder {
	builder.getTokenOutputBuilder.Scope = scope
	return builder
}
