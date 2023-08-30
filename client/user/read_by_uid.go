package user

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadUserOutput, error) {

	readReq := NewReadByUidRequest(ctx, client, readInp.Uid)
	var userDetails UserDetails
	if readErr := readReq.Send(&userDetails); readErr != nil {
		return nil, readErr
	}

	return &userDetails, nil
}
