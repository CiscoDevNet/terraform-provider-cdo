package sdc

import (
	"crypto/rand"
	"crypto/rsa"

	internalRsa "github.com/CiscoDevnet/go-client/internal/crypto/rsa"
)

type sdcReadOutputBuilder struct {
	readOutput ReadOutput
}

func NewSdcOutputBuilder() *sdcReadOutputBuilder {
	return &sdcReadOutputBuilder{}
}

func (builder *sdcReadOutputBuilder) Build() ReadOutput {
	return builder.readOutput
}

func (builder *sdcReadOutputBuilder) AsDefaultCloudConnector() *sdcReadOutputBuilder {
	builder.readOutput.DefaultSdc = true
	builder.readOutput.Cdg = true

	return builder
}

func (builder *sdcReadOutputBuilder) AsOnPremConnector() *sdcReadOutputBuilder {
	builder.readOutput.Cdg = false
	builder.readOutput.PublicKey = PublicKey{
		EncodedKey: mustGenerateBase64PublicKey(),
		Version:    164,
		KeyId:      "01010101-0101-0101-0101-010101010101",
	}

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

func mustGenerateBase64PublicKey() string {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		panic(err)
	}

	return internalRsa.MustBase64PublicKeyFromRsaKey(key)
}
