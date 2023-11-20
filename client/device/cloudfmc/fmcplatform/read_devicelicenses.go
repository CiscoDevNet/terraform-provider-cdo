package fmcplatform

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/devicelicense"
)

type ReadDeviceLicensesInput struct {
	FmcHost string
}

func NewReadDeviceLicensesInput(fmcHost string) ReadDeviceLicensesInput {
	return ReadDeviceLicensesInput{
		FmcHost: fmcHost,
	}
}

type ReadDeviceLicensesOutput = devicelicense.Item

var NewReadDeviceLicensesOutputBuilder = devicelicense.NewItemBuilder

func ReadDeviceLicenses(ctx context.Context, client http.Client, readInp ReadDeviceLicensesInput) (*ReadDeviceLicensesOutput, error) {

	client.Logger.Println("reading FMC device licenses")

	readUrl := url.ReadFmcDeviceLicenses(client.BaseUrl())
	req := client.NewGet(ctx, readUrl)
	req.Header.Set("Fmc-Hostname", readInp.FmcHost)

	var outp devicelicense.DeviceLicense
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	if len(outp.Items) != 1 {
		return nil, fmt.Errorf("failed to get device license")
	}

	return &outp.Items[0], nil
}
