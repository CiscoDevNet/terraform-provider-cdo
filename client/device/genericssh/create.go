package genericssh

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/goutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type CreateInput struct {
	Name          string
	ConnectorUid  string
	ConnectorType string
	SocketAddress string
	Tags          tags.Type
}

type CreateOutput = device.CreateOutput

func NewCreateInput(name, connectorUid, socketAddress string, tags tags.Type) CreateInput {
	return CreateInput{
		Name:          name,
		ConnectorUid:  connectorUid,
		SocketAddress: socketAddress,
		Tags:          tags,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating generic ssh")

	deviceInput := device.NewCreateInputBuilder().
		Name(createInp.Name).
		DeviceType(devicetype.GenericSSH).
		ConnectorUid(createInp.ConnectorUid).
		ConnectorType(createInp.ConnectorType).
		SocketAddress(createInp.SocketAddress).
		Model(false).
		IgnoreCertificate(goutil.NewBoolPointer(false)).
		Metadata(nil).
		Tags(createInp.Tags).
		Build()

	outp, err := device.Create(ctx, client, deviceInput)
	if err != nil {
		return nil, err
	}

	return outp, nil
}
