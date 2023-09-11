package cloudftdonboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type ReadInput struct {
}

func NewReadInput() ReadInput {
	return ReadInput{}
}

type ReadOutput struct {
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {
	return &ReadOutput{}, nil
}
