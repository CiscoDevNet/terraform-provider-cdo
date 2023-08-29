package user

import (
	"context"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

func Create(ctx context.Context, client http.Client, createInp CreateUserInput) (*UserDetails, error) {
	client.Logger.Println(fmt.Sprintf("Creating user %s", createInp.Username))
	req := NewCreateRequest(ctx, client, createInp)

	var outp UserTenantAssociation
	if err := req.SendFormUrlEncoded(&outp); err != nil {
		return nil, err
	}

	// the user creation endpoint annoyingly returns an association, so we need to make a read request to get the actual user
	readReq := NewReadByUidRequest(ctx, client, outp.Source.Uid)
	var readOutp UserDetails
	if readErr := readReq.Send(&readOutp); readErr != nil {
		return nil, readErr
	}

	return &readOutp, nil
}
