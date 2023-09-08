package fmcconfig

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

type CreateDeviceRecordInputBuilder struct {
	createDeviceRecordInput *CreateDeviceRecordInput
}

func NewCreateDeviceRecordInputBuilder() *CreateDeviceRecordInputBuilder {
	createDeviceRecordInput := &CreateDeviceRecordInput{}
	b := &CreateDeviceRecordInputBuilder{createDeviceRecordInput: createDeviceRecordInput}
	return b
}

func (b *CreateDeviceRecordInputBuilder) FmcDomainUid(fmcDomainUid string) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.FmcDomainUid = fmcDomainUid
	return b
}

func (b *CreateDeviceRecordInputBuilder) SystemApiToken(systemApiToken string) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.SystemApiToken = systemApiToken
	return b
}

func (b *CreateDeviceRecordInputBuilder) Name(name string) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.Name = name
	return b
}

func (b *CreateDeviceRecordInputBuilder) NatId(natId string) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.NatId = natId
	return b
}

func (b *CreateDeviceRecordInputBuilder) RegKey(regKey string) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.RegKey = regKey
	return b
}

func (b *CreateDeviceRecordInputBuilder) PerformanceTier(performanceTier *tier.Type) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.PerformanceTier = performanceTier
	return b
}

func (b *CreateDeviceRecordInputBuilder) LicenseCaps(licenseCaps []license.Type) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.LicenseCaps = licenseCaps
	return b
}

func (b *CreateDeviceRecordInputBuilder) AccessPolicyUid(accessPolicyUid string) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.AccessPolicyUid = accessPolicyUid
	return b
}

func (b *CreateDeviceRecordInputBuilder) Type(_type string) *CreateDeviceRecordInputBuilder {
	b.createDeviceRecordInput.Type = _type
	return b
}

func (b *CreateDeviceRecordInputBuilder) Build() CreateDeviceRecordInput {
	return *b.createDeviceRecordInput
}
