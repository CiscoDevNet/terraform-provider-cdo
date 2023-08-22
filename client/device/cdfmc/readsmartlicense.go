package cdfmc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/smartlicense"
)

type ReadInput struct{}

type ReadOutput = smartlicense.SmartLicense

func ReadSmartLicense(ctx context.Context, client http.Client, _inp ReadInput) (*ReadOutput, error) {
	readUrl := url.ReadSmartLicense(client.BaseUrl())
	
	req := client.NewGet(ctx, readUrl)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
