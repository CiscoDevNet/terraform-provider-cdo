package cloudftd

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
)

type UpdateInputBuilder struct {
	updateInput *UpdateInput
}

func NewUpdateInputBuilder() *UpdateInputBuilder {
	updateInput := &UpdateInput{}
	b := &UpdateInputBuilder{updateInput: updateInput}
	return b
}

func (b *UpdateInputBuilder) Uid(uid string) *UpdateInputBuilder {
	b.updateInput.Uid = uid
	return b
}

func (b *UpdateInputBuilder) Name(name string) *UpdateInputBuilder {
	b.updateInput.Name = name
	return b
}

func (b *UpdateInputBuilder) Tags(tags tags.Type) *UpdateInputBuilder {
	b.updateInput.Tags = tags
	return b
}

func (b *UpdateInputBuilder) Licenses(licenses []license.Type) *UpdateInputBuilder {
	b.updateInput.Licenses = licenses
	return b
}

func (b *UpdateInputBuilder) Build() UpdateInput {
	return *b.updateInput
}
