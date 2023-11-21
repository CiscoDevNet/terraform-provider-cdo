package devicelicense

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
)

type ItemBuilder struct {
	item *Item
}

func NewItemBuilder() *ItemBuilder {
	item := &Item{}
	b := &ItemBuilder{item: item}
	return b
}

func (b *ItemBuilder) Id(id string) *ItemBuilder {
	b.item.Id = id
	return b
}

func (b *ItemBuilder) Type(type_ string) *ItemBuilder {
	b.item.Type = type_
	return b
}

func (b *ItemBuilder) LicenseTypes(licenseTypes []license.Type) *ItemBuilder {
	b.item.LicenseTypes = licenseTypes
	return b
}

func (b *ItemBuilder) PerformanceTier(performanceTier tier.Type) *ItemBuilder {
	b.item.PerformanceTier = performanceTier
	return b
}

func (b *ItemBuilder) Links(links Links) *ItemBuilder {
	b.item.Links = links
	return b
}

func (b *ItemBuilder) Build() Item {
	return *b.item
}
