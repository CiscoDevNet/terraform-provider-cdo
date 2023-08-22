package device

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type ReadAllByTypeInput struct {
	DeviceType string `json:"deviceType"`
}

type ReadAllOutput = []ReadOutput

func NewReadAllByTypeInput(deviceType string) ReadAllByTypeInput {
	return ReadAllByTypeInput{
		DeviceType: deviceType,
	}
}

func ReadAllByType(ctx context.Context, client http.Client, readInp ReadAllByTypeInput) (*ReadAllOutput, error) {

	client.Logger.Println("reading all Devices by device type")

	readAllUrl := url.ReadAllDevicesByType(client.BaseUrl(), readInp.DeviceType)

	req := client.NewGet(ctx, readAllUrl)

	var outp []ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
