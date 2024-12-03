package tenants

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func AddExistingTenantUsingApiToken(ctx context.Context, client http.Client, addInp MspAddExistingTenantInput) (*MspTenantOutput, *CreateError) {
	client.Logger.Println("Creating tenant for CDO")
	addUrl := url.AddExistingTenantToMspManagedTenant(client.BaseUrl())

	req := client.NewPost(ctx, addUrl, addInp)

	var createOutp MspManagedTenantStatusInfo
	if err := req.Send(&createOutp); err != nil {
		return nil, &CreateError{Err: err}
	}

	return &createOutp.MspManagedTenant, nil
}

type StatusInfo struct {
	status string
}
