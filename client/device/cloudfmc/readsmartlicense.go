package cloudfmc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/smartlicense"
)

type ReadSmartLicenseInput struct{}

func NewReadSmartLicenseInput() ReadSmartLicenseInput {
	return ReadSmartLicenseInput{}
}

type ReadSmartLicenseOutput = smartlicense.SmartLicense

func ReadSmartLicense(ctx context.Context, client http.Client, _inp ReadSmartLicenseInput) (*ReadSmartLicenseOutput, error) {
	readUrl := url.ReadSmartLicense(client.BaseUrl())

	req := client.NewGet(ctx, readUrl)

	var outp ReadSmartLicenseOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
