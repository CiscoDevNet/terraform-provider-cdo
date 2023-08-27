package cloudfmc

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

type ReadOutput = device.ReadOutput

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading cloud FMC")

	req := device.ReadAllByTypeRequest(ctx, client, device.NewReadAllByTypeInput(devicetype.CloudFmc))
	var cloudFmcDevices []ReadOutput
	if err := req.Send(&cloudFmcDevices); err != nil {
		return nil, err
	}

	if len(cloudFmcDevices) == 0 {
		return nil, fmt.Errorf("firewall management center (FMC) not found")
	}

	if len(cloudFmcDevices) > 1 {
		return nil, fmt.Errorf("more than one firewall management center (FMC) found, please report this issue at: %s", cdo.TerraformProviderCDOIssuesUrl)
	}

	return &cloudFmcDevices[0], nil
}
