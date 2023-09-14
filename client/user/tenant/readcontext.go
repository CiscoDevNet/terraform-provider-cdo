package tenant

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/tenant"
)

type ReadContextInput struct {
}

func NewReadContextInput() ReadContextInput {
	return ReadContextInput{}
}

type ReadContextOutput = tenant.Context

func ReadContext(ctx context.Context, client http.Client, _ ReadContextInput) (*ReadContextOutput, error) {

	readUrl := url.ReadTenantContext(client.BaseUrl())

	req := client.NewGet(ctx, readUrl)

	var readOutp []ReadContextOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	if len(readOutp) < 1 {
		return nil, fmt.Errorf("tenant context not found")
	}

	// TODO: Question: is this a valid case?
	if len(readOutp) > 1 {
		return nil, fmt.Errorf("more than one tenant context found")
	}

	return &readOutp[0], nil
}
