package device

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/goutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
)

// CreateInput builder pattern code
type CreateInputBuilder struct {
	createInput *CreateInput
}

func NewCreateInputBuilder() *CreateInputBuilder {
	createInput := &CreateInput{}
	b := &CreateInputBuilder{createInput: createInput}
	return b
}

func (b *CreateInputBuilder) Name(name string) *CreateInputBuilder {
	b.createInput.Name = name
	return b
}

func (b *CreateInputBuilder) DeviceType(deviceType devicetype.Type) *CreateInputBuilder {
	b.createInput.DeviceType = deviceType
	return b
}

func (b *CreateInputBuilder) Model(model bool) *CreateInputBuilder {
	b.createInput.Model = model
	return b
}

func (b *CreateInputBuilder) ConnectorUid(connectorUid string) *CreateInputBuilder {
	b.createInput.ConnectorUid = connectorUid
	return b
}

func (b *CreateInputBuilder) ConnectorType(connectorType string) *CreateInputBuilder {
	b.createInput.ConnectorType = connectorType
	return b
}

func (b *CreateInputBuilder) SocketAddress(socketAddress string) *CreateInputBuilder {
	b.createInput.SocketAddress = socketAddress
	return b
}

func (b *CreateInputBuilder) IgnoreCertificate(ignoreCertificate *bool) *CreateInputBuilder {
	b.createInput.IgnoreCertificate = ignoreCertificate
	return b
}

func (b *CreateInputBuilder) Metadata(metadata interface{}) *CreateInputBuilder {
	b.createInput.Metadata = goutil.AsPointer(metadata)
	return b
}

func (b *CreateInputBuilder) Tags(tags tags.Type) *CreateInputBuilder {
	b.createInput.Tags = tags
	return b
}

func (b *CreateInputBuilder) EnableOobDetection(enableOobDetection *bool) *CreateInputBuilder {
	b.createInput.EnableOobDetection = enableOobDetection
	return b
}

func (b *CreateInputBuilder) Build() CreateInput {
	return *b.createInput
}
