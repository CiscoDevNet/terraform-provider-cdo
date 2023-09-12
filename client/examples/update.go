package examples

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type UpdateInput struct {
}

func NewUpdateInput() UpdateInput {
	return UpdateInput{}
}

type UpdateOutput struct {
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	// TODO

	return nil, nil
}
