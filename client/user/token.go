package user

import (
	"context"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type GetTokenInput struct{}

type GetTokenOutput struct {
	TenantUid    string `json:"tenantUid"`
	TenantName   string `json:"tenantName"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

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
