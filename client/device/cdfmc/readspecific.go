package cdfmc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type ReadSpecificInput struct {
	FmcId string
}

type ReadSpecificOutput struct {
	Uid       string `json:"uid"`
	DomainUid string `json:"domainUid"`
	State     string `json:"state"`
	Status    string `json:"status"`
}

func NewReadSpecificInput(fmcId string) ReadSpecificInput {
	return ReadSpecificInput{
		FmcId: fmcId,
	}
}

func ReadSpecific(ctx context.Context, client http.Client, inp ReadSpecificInput) (*ReadSpecificOutput, error) {

	req := device.NewReadSpecificRequest(ctx, client, *device.NewReadSpecificInput(inp.FmcId))

	var readSpecificOutp ReadSpecificOutput
	if err := req.Send(&readSpecificOutp); err != nil {
		return nil, err
	}

	return &readSpecificOutp, nil
}
