package cloudfmc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/accesspolicies"
	"strconv"
)

type ReadAccessPoliciesInput struct {
	FmcHostname string
	DomainUid   string
	Limit       int
}

func NewReadAccessPoliciesInput(fmcHostname, domainUid string, limit int) ReadAccessPoliciesInput {
	return ReadAccessPoliciesInput{
		FmcHostname: fmcHostname,
		DomainUid:   domainUid,
		Limit:       limit,
	}
}

type ReadAccessPoliciesOutput = accesspolicies.AccessPolicies

func ReadAccessPolicies(ctx context.Context, client http.Client, inp ReadAccessPoliciesInput) (*ReadAccessPoliciesOutput, error) {

	readUrl := url.ReadAccessPolicies(client.BaseUrl(), inp.DomainUid)

	req := client.NewGet(ctx, readUrl)
	req.Header.Add("Fmc-Hostname", inp.FmcHostname) // required, otherwise 500 internal server error
	req.QueryParams.Add("limit", strconv.Itoa(inp.Limit))

	var outp ReadAccessPoliciesOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
