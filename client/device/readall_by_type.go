package device

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
)

type ReadAllByTypeInput struct {
	DeviceType devicetype.Type `json:"deviceType"`
}

type ReadAllByTypeOutput = []ReadOutput

func NewReadAllByTypeInput(deviceType devicetype.Type) ReadAllByTypeInput {
	return ReadAllByTypeInput{
		DeviceType: deviceType,
	}
}

func ReadAllByType(ctx context.Context, client http.Client, readInp ReadAllByTypeInput) (*ReadAllByTypeOutput, error) {

	client.Logger.Println("reading all Devices by device type")

	readAllUrl := url.ReadAllDevicesByType(client.BaseUrl(), readInp.DeviceType)

	req := client.NewGet(ctx, readAllUrl)

	var outp []ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
