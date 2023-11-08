package sec

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type UpdateInput struct {
}

type UpdateOutput struct {
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	// intentional empty, I don't know what update we can do on an SEC

	return &UpdateOutput{}, nil
}
