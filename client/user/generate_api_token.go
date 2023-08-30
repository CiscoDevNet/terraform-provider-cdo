package user

import (
	"context"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

func GenerateApiToken(ctx context.Context, client http.Client, generateApiTokenInp GenerateApiTokenInput) (*ApiTokenResponse, error) {
	client.Logger.Println(fmt.Sprintf("Generating API token for user %s", generateApiTokenInp.Name))
	req := NewGenerateApiTokenRequest(ctx, client, generateApiTokenInp)

	var outp ApiTokenResponse
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
