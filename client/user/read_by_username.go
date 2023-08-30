package user

import (
	"context"
	"errors"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

func ReadByUsername(ctx context.Context, client http.Client, readInp ReadByUsernameInput) (*ReadUserOutput, error) {

	readReq := NewReadByUsernameRequest(ctx, client, readInp.Name)
	var userDetails []UserDetails
	if readErr := readReq.Send(&userDetails); readErr != nil {
		return nil, readErr
	}
	if len(userDetails) != 1 {
		return nil, errors.New("User not found")
	}

	return &userDetails[0], nil
}
