package user

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func Delete(ctx context.Context, client http.Client, deleteInp DeleteUserInput) (*DeleteUserOutput, error) {

	url := url.UserByUid(client.BaseUrl(), deleteInp.Uid)

	req := client.NewDelete(ctx, url)

	var deleteOutp DeleteUserOutput
	if err := req.Send(&deleteOutp); err != nil {
		return nil, err
	}

	return &deleteOutp, nil
}
