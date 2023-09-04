package cloudftd

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

type MetadataBuilder struct {
	metadata *Metadata
}

func NewMetadataBuilder() *MetadataBuilder {
	metadata := &Metadata{}
	b := &MetadataBuilder{metadata: metadata}
	return b
}

func (b *MetadataBuilder) AccessPolicyName(accessPolicyName string) *MetadataBuilder {
	b.metadata.AccessPolicyName = accessPolicyName
	return b
}

func (b *MetadataBuilder) AccessPolicyUuid(accessPolicyUuid string) *MetadataBuilder {
	b.metadata.AccessPolicyUid = accessPolicyUuid
	return b
}

func (b *MetadataBuilder) CloudManagerDomain(cloudManagerDomain string) *MetadataBuilder {
	b.metadata.CloudManagerDomain = cloudManagerDomain
	return b
}

func (b *MetadataBuilder) GeneratedCommand(generatedCommand string) *MetadataBuilder {
	b.metadata.GeneratedCommand = generatedCommand
	return b
}

func (b *MetadataBuilder) LicenseCaps(licenseCaps []license.Type) *MetadataBuilder {
	b.metadata.LicenseCaps = licenseCaps
	return b
}

func (b *MetadataBuilder) NatID(natID string) *MetadataBuilder {
	b.metadata.NatID = natID
	return b
}

func (b *MetadataBuilder) PerformanceTier(performanceTier *tier.Type) *MetadataBuilder {
	b.metadata.PerformanceTier = performanceTier
	return b
}

func (b *MetadataBuilder) RegKey(regKey string) *MetadataBuilder {
	b.metadata.RegKey = regKey
	return b
}

func (b *MetadataBuilder) Build() Metadata {
	return *b.metadata
}
