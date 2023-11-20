package fmcplatform

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/devicelicense"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
)

type UpdateDeviceLicensesInput struct {
	FmcHost      string
	LicenseTypes []license.Type
}

type UpdateDeviceLicensesOutput = devicelicense.Item

type updateRequestBody struct {
	Type         string         `json:"type"`
	Id           string         `json:"id"`
	LicenseTypes []license.Type `json:"licenseTypes"` // must be FMC license types
}

func UpdateDeviceLicenses(ctx context.Context, client http.Client, updateInp UpdateDeviceLicensesInput) (*UpdateDeviceLicensesOutput, error) {

	client.Logger.Println("updating FMC device licenses")

	readOutput, err := ReadDeviceLicenses(ctx, client, NewReadDeviceLicensesInputBuilder().FmcHost(updateInp.FmcHost).Build())
	if err != nil {
		return nil, err
	}

	updateUrl := url.UpdateFmcDeviceLicenses(client.BaseUrl(), readOutput.Id)
	updateBody := updateRequestBody{
		Type:         "DeviceLicense",
		Id:           readOutput.Id,
		LicenseTypes: license.LicensesToFmcLicenses(updateInp.LicenseTypes),
	}
	req := client.NewPut(ctx, updateUrl, updateBody)
	req.Header.Set("Fmc-Hostname", updateInp.FmcHost)

	var outp UpdateDeviceLicensesOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
