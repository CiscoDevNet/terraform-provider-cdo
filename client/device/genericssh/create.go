package genericssh

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type CreateInput struct {
	Name          string
	ConnectorUid  string
	ConnectorType string
	SocketAddress string
}

type CreateOutput = device.CreateOutput

func NewCreateInput(name, connectorUid, socketAddress string) CreateInput {
	return CreateInput{
		Name:          name,
		ConnectorUid:  connectorUid,
		SocketAddress: socketAddress,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating generic ssh")

	deviceInput := device.NewCreateRequestInput(createInp.Name, "GENERIC_SSH", createInp.ConnectorUid, createInp.ConnectorType, createInp.SocketAddress, false, false, nil)
	outp, err := device.Create(ctx, client, *deviceInput)
	if err != nil {
		return nil, err
	}

	return outp, nil
}
