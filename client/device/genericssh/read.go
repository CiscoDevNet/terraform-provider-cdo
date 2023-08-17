package genericssh

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type ReadInput struct {
	Uid string `json:"-"`
}

type ReadOutput = device.ReadOutput

func NewReadInput(uid string) *UpdateInput {
	return &UpdateInput{
		Uid: uid,
	}
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("updating generic ssh")

	readOutp, err := device.Read(ctx, client, *device.NewReadInput(readInp.Uid))
	if err != nil {
		return nil, err
	}

	return readOutp, nil
}
