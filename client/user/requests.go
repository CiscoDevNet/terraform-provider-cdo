package user

import (
	"context"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type updateRequestBody struct {
	UserRoles []string `json:"roles"`
}

func NewCreateRequest(ctx context.Context, client http.Client, createInp CreateUserInput) *http.Request {
	url := url.CreateUser(client.BaseUrl(), createInp.Username)
	body := fmt.Sprintf("roles=%s&isApiOnlyUser=%t", createInp.UserRoles, createInp.ApiOnlyUser)
	return client.NewPost(ctx, url, body)
}

func NewGenerateApiTokenRequest(ctx context.Context, client http.Client, generateApiTokenInp GenerateApiTokenInput) *http.Request {
	url := url.GenerateApiToken(client.BaseUrl(), generateApiTokenInp.Name)
	var body = struct{}{}
	return client.NewPost(ctx, url, body)
}

func NewRevokeOauthTokenRequest(ctx context.Context, client http.Client, revokeApiTokenInp RevokeOAuthTokenInput) *http.Request {
	url := url.RevokeApiToken(client.BaseUrl(), revokeApiTokenInp.ApiTokenId)
	var body = struct{}{}
	return client.NewPost(ctx, url, body)
}

func NewReadByUidRequest(ctx context.Context, client http.Client, uid string) *http.Request {
	url := url.UserByUid(client.BaseUrl(), uid)
	return client.NewGet(ctx, url)
}

func NewReadByUsernameRequest(ctx context.Context, client http.Client, username string) *http.Request {
	url := url.ReadUserByUsername(client.BaseUrl(), username)
	return client.NewGet(ctx, url)
}

func NewUpdateRequest(ctx context.Context, client http.Client, updateInp UpdateUserInput) *http.Request {
	url := url.UserByUid(client.BaseUrl(), updateInp.Uid)
	body := updateRequestBody{
		UserRoles: updateInp.UserRoles,
	}

	return client.NewPut(ctx, url, body)
}
