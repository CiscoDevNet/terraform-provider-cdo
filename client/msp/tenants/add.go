package tenants

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func AddExistingTenantUsingApiToken(ctx context.Context, client http.Client, addInp MspAddExistingTenantInput) (*MspTenantOutput, *CreateError) {
	client.Logger.Println("Adding existing tenant to MSp portal...")
	addUrl := url.AddExistingTenantToMspManagedTenant(client.BaseUrl())

	req := client.NewPost(ctx, addUrl, addInp)

	var createOutp MspManagedTenantStatusInfo
	if err := req.Send(&createOutp); err != nil {
		return nil, &CreateError{Err: err}
	}

	client.Logger.Printf("Added existing tenant %s to MSP portal using API token...", createOutp.MspManagedTenant.Name)

	return &createOutp.MspManagedTenant, nil
}
