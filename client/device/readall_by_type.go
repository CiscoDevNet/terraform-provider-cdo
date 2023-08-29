package device

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
)

type ReadAllByTypeInput struct {
	DeviceType devicetype.Type
}

type ReadAllByTypeOutput = []ReadOutput

func NewReadAllByTypeInput(deviceType devicetype.Type) ReadAllByTypeInput {
	return ReadAllByTypeInput{
		DeviceType: deviceType,
	}
}

func ReadAllByTypeRequest(ctx context.Context, client http.Client, readInp ReadAllByTypeInput) *http.Request {
	readAllUrl := url.ReadAllDevicesByType(client.BaseUrl())

	req := client.NewGet(ctx, readAllUrl)

	req.QueryParams.Add("q", fmt.Sprintf("deviceType:%s", readInp.DeviceType))

	return req
}

func ReadAllByType(ctx context.Context, client http.Client, readInp ReadAllByTypeInput) (*ReadAllByTypeOutput, error) {

	client.Logger.Println("reading all Devices by device type")

	req := ReadAllByTypeRequest(ctx, client, readInp)

	var outp []ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
