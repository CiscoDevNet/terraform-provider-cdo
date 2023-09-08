package user

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth"
)

type GetTokenInfoInput struct {
}

func NewGetTokenInfoInput() GetTokenInfoInput {
	return GetTokenInfoInput{}
}

type GetTokenInfoOutput = auth.Info

func GetTokenInfo(ctx context.Context, client http.Client, getTokenInfoInp GetTokenInfoInput) (*GetTokenInfoOutput, error) {

	readUrl := url.ReadTokenInfo(client.BaseUrl())

	req := client.NewGet(ctx, readUrl)

	var readOutp GetTokenInfoOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}
