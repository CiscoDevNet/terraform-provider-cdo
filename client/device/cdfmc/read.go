package cdfmc

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type ReadInput struct {
	Uid string `json:"uid"`
}

type ReadOutput = device.ReadOutput

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	cdfmcDevices, err := device.ReadAllByType(ctx, client, device.NewReadAllByTypeInput("FMCE"))
	if err != nil {
		return nil, err
	}

	if len(*cdfmcDevices) == 0 {
		return nil, fmt.Errorf("firewall management center (FMC) not found")
	}

	if len(*cdfmcDevices) > 1 {
		return nil, fmt.Errorf("more than one firewall management center (FMC) found, please report this issue at: %s", cdo.TerraformProviderCDOIssuesUrl)
	}

	return &(*cdfmcDevices)[0], nil
}
