package genericssh

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"

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

	deviceInput := device.NewCreateRequestInput(
		createInp.Name,
		"GENERIC_SSH",
		createInp.ConnectorUid,
		createInp.ConnectorType,
		createInp.SocketAddress,
		false,
		false,
		nil,
		createInp.Tags,
	)
	outp, err := device.Create(ctx, client, *deviceInput)
	if err != nil {
		return nil, err
	}

	return outp, nil
}
