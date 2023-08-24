package device

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
	"time"
)

type readOutputBuilder struct {
	readOutput ReadOutput
}

func NewReadOutputBuilder() *readOutputBuilder {
	return &readOutputBuilder{
		readOutput: CreateOutput{
			CreatedDate:     time.Now().Unix(),
			LastUpdatedDate: time.Now().Unix(),
		},
	}
}

func (builder *readOutputBuilder) Build() ReadOutput {
	return builder.readOutput
}

func (builder *readOutputBuilder) AsAsa() *readOutputBuilder {
	builder.readOutput.DeviceType = devicetype.Asa
	return builder
}

func (builder *readOutputBuilder) AsIos() *readOutputBuilder {
	builder.readOutput.DeviceType = "IOS"
	return builder
}

func (builder *readOutputBuilder) AsCdfmc() *readOutputBuilder {
	builder.readOutput.DeviceType = "FMCE"
	return builder
}

func (builder *readOutputBuilder) WithUid(uid string) *readOutputBuilder {
	builder.readOutput.Uid = uid

	return builder
}

func (builder *readOutputBuilder) WithName(name string) *readOutputBuilder {
	builder.readOutput.Name = name

	return builder
}

func (builder *readOutputBuilder) WithLocation(host string, port uint) *readOutputBuilder {
	builder.readOutput.Host = host
	builder.readOutput.Port = fmt.Sprint(port)
	builder.readOutput.SocketAddress = fmt.Sprintf("%s:%d", host, port)

	return builder
}

func (builder *readOutputBuilder) WithCreatedDate(date time.Time) *readOutputBuilder {
	builder.readOutput.CreatedDate = date.Unix()

	return builder
}

func (builder *readOutputBuilder) WithLastUpdatedDate(date time.Time) *readOutputBuilder {
	builder.readOutput.LastUpdatedDate = date.Unix()

	return builder
}

func (builder *readOutputBuilder) OnboardedUsingOnPremConnector(connectorUid string) *readOutputBuilder {
	builder.readOutput.ConnectorType = "SDC"
	builder.readOutput.ConnectorUid = connectorUid

	return builder
}

func (builder *readOutputBuilder) OnboardedUsingCloudConnector(cdgUid string) *readOutputBuilder {
	builder.readOutput.ConnectorType = "CDG"
	builder.readOutput.ConnectorUid = cdgUid

	return builder
}
