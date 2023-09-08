package user

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth"
)

type CreateSystemTokenInput struct {
	Scope string
}

func NewCreateSystemTokenInput(scope string) CreateSystemTokenInput {
	return CreateSystemTokenInput{
		Scope: scope,
	}
}

type CreateSystemTokenOutput = auth.Token

// CreateSystemToken returns a system token that is valid for 24 hrs
func CreateSystemToken(ctx context.Context, client http.Client, getSystemTokenInp CreateSystemTokenInput) (*CreateSystemTokenOutput, error) {

	createUrl := url.CreateSystemToken(client.BaseUrl(), getSystemTokenInp.Scope)

	createReq := client.NewPost(ctx, createUrl, nil)
	var createOutp CreateSystemTokenOutput
	if err := createReq.Send(&createOutp); err != nil {
		return nil, err
	}

	return &createOutp, nil
}
