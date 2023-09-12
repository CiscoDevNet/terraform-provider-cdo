package user

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type GetTokenInput struct{}

type GetTokenOutput = auth.Token

func NewGetTokenInput() GetTokenInput {
	return GetTokenInput{}
}

func GetToken(ctx context.Context, client http.Client, getTokenInp GetTokenInput) (GetTokenOutput, error) {

	url := url.UserToken(client.BaseUrl())

	req := client.NewPost(ctx, url, nil)

	var getTokenOutp GetTokenOutput
	if err := req.Send(&getTokenOutp); err != nil {
		return GetTokenOutput{}, err
	}

	return getTokenOutp, nil
}
