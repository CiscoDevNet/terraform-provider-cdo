package user

import (
	"context"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

func Update(ctx context.Context, client http.Client, updateInp UpdateUserInput) (*UserDetails, error) {
	client.Logger.Println(fmt.Sprintf("Creating user %s", updateInp.Uid))

	req := NewUpdateRequest(ctx, client, updateInp)
	var outp UserTenantAssociation

	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	// the user update endpoint annoyingly returns an association, so we need to make a read request to get the actual user
	readReq := NewReadByUidRequest(ctx, client, outp.Source.Uid)
	var readOutp UserDetails
	if readErr := readReq.Send(&readOutp); readErr != nil {
		return nil, readErr
	}

	return &readOutp, nil
}
