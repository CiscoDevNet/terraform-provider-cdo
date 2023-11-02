package connector

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/crypto"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
)

type sdcReadOutputBuilder struct {
	readOutput ReadOutput
}

func NewConnectorOutputBuilder() *sdcReadOutputBuilder {
	return &sdcReadOutputBuilder{}
}

func (builder *sdcReadOutputBuilder) Build() ReadOutput {
	return builder.readOutput
}

func (builder *sdcReadOutputBuilder) AsDefaultCloudConnector() *sdcReadOutputBuilder {
	builder.readOutput.DefaultConnector = true
	builder.readOutput.Cdg = true

	return builder
}

func (builder *sdcReadOutputBuilder) AsOnPremConnector() *sdcReadOutputBuilder {
	builder.readOutput.Cdg = false
	builder.readOutput.PublicKey = model.NewPublicKey(mustGenerateBase64PublicKey(), 164, "01010101-0101-0101-0101-010101010101")

	return builder
}

func (builder *sdcReadOutputBuilder) WithUid(uid string) *sdcReadOutputBuilder {
	builder.readOutput.Uid = uid

	return builder
}

func (builder *sdcReadOutputBuilder) WithName(name string) *sdcReadOutputBuilder {
	builder.readOutput.Name = name

	return builder
}

func (builder *sdcReadOutputBuilder) WithTenantUid(tenantUid string) *sdcReadOutputBuilder {
	builder.readOutput.TenantUid = tenantUid

	return builder
}

func (builder *sdcReadOutputBuilder) WithCommunicationReady(ready bool) *sdcReadOutputBuilder {
	builder.readOutput.IsCommunicationQueueReady = ready

	return builder
}

func mustGenerateBase64PublicKey() string {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	return crypto.MustBase64PublicKeyFromRsaKey(key)
}
