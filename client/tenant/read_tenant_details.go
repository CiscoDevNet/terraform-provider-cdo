package tenant

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func ReadTenantDetails(ctx context.Context, client http.Client) (*ReadTenantDetailsOutput, error) {
	client.Logger.Println("Get tenant details for currently connected client")

	req := NewReadTenantDetailsRequest(ctx, client)

	var outp ReadTenantDetailsOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}

func NewReadTenantDetailsRequest(ctx context.Context, client http.Client) *http.Request {
	url := url.ReadTenantDetails(client.BaseUrl())
	return client.NewGet(ctx, url)
}
