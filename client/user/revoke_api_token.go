package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

func RevokeApiToken(ctx context.Context, client http.Client, revokeInp RevokeApiTokenInput) (*RevokeApiTokenOutput, error) {
	client.Logger.Println(fmt.Sprintf("Revoking API token for %s", revokeInp.Name))

	// 1. Find the user
	readReq := NewReadByUsernameRequest(ctx, client, revokeInp.Name)
	var userDetails []UserDetails
	if readErr := readReq.Send(&userDetails); readErr != nil {
		return nil, readErr
	}
	if len(userDetails) != 1 {
		return nil, errors.New("User not found")
	}

	// 2. Revoke the API token by ID for the user
	revokeApiTokenInput := NewRevokeOAuthTokenInput(userDetails[0].ApiTokenId)
	revokeReq := NewRevokeOauthTokenRequest(ctx, client, *revokeApiTokenInput)
	var revokeOutput RevokeApiTokenOutput
	if revokeErr := revokeReq.Send(&revokeOutput); revokeErr != nil {
		return nil, revokeErr
	}

	return &revokeOutput, nil
}
