package cdfmc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/accesspolicies"
)

type ReadAccessPoliciesInput struct {
	DomainUid string
	Limit     int
}

func NewReadAccessPoliciesInput(domainUid string, limit int) ReadAccessPoliciesInput {
	return ReadAccessPoliciesInput{
		DomainUid: domainUid,
		Limit:     limit,
	}
}

type ReadAccessPoliciesOutput = accesspolicies.AccessPolicies

func ReadAccessPolicies(ctx context.Context, client http.Client, inp ReadAccessPoliciesInput) (*ReadAccessPoliciesOutput, error) {

	readUrl := url.ReadAccessPolicies(client.BaseUrl(), inp.DomainUid, inp.Limit)

	req := client.NewGet(ctx, readUrl)

	var outp ReadAccessPoliciesOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
