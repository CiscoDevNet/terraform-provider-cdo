package cdfmc

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type ReadInput struct {
}

func NewReadInput() ReadInput {
	return ReadInput{}
}

type ReadOutput struct {
	device.ReadOutput
	Host string `json:"host"`
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading cdFMC")

	req := device.ReadAllByTypeRequest(ctx, client, device.NewReadAllByTypeInput(devicetype.Cdfmc))
	var cdfmcDevices []ReadOutput
	if err := req.Send(&cdfmcDevices); err != nil {
		return nil, err
	}

	if len(cdfmcDevices) == 0 {
		return nil, fmt.Errorf("firewall management center (FMC) not found")
	}

	if len(cdfmcDevices) > 1 {
		return nil, fmt.Errorf("more than one firewall management center (FMC) found, please report this issue at: %s", cdo.TerraformProviderCDOIssuesUrl)
	}

	return &cdfmcDevices[0], nil
}
