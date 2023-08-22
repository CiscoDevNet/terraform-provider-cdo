package device

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type ReadByNameAndTypeInput struct {
	Name       string `json:"name"`
	DeviceType string `json:"deviceType"`
}

func NewReadByNameAndTypeInput(name, deviceType string) ReadByNameAndTypeInput {
	return ReadByNameAndTypeInput{
		Name:       name,
		DeviceType: deviceType,
	}
}

func ReadByNameAndType(ctx context.Context, client http.Client, readInp ReadByNameAndTypeInput) (*ReadOutput, error) {

	client.Logger.Println("reading Device by name and device type")

	readUrl := url.ReadDeviceByNameAndDeviceType(client.BaseUrl(), readInp.Name, readInp.DeviceType)

	req := client.NewGet(ctx, readUrl)

	var arrayOutp []ReadOutput
	if err := req.Send(&arrayOutp); err != nil {
		return nil, err
	}

	if len(arrayOutp) == 0 {
		return nil, fmt.Errorf("no Device by name %s and device type %s found", readInp.Name, readInp.DeviceType)
	}

	if len(arrayOutp) > 1 {
		return nil, fmt.Errorf("multiple devices found with the name: %s and device type: %s", readInp.Name, readInp.DeviceType)
	}

	outp := arrayOutp[0]
	return &outp, nil
}
