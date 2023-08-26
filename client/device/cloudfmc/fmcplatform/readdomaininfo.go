package fmcplatform

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcdomain"
)

// TODO: use this to fetch
//
//	curl --request GET \
//	 --url https://<FMC_HOST>/api/fmc_platform/v1/info/domain \
//	 --header 'Authorization: Bearer <CDO TOKEN>'

type ReadDomainInfoInput struct {
	FmcHost string
}

func NewReadDomainInfo(fmcHost string) ReadDomainInfoInput {
	return ReadDomainInfoInput{
		FmcHost: fmcHost,
	}
}

type ReadDomainInfoOutput = fmcdomain.Info

func ReadFmcDomainInfo(ctx context.Context, client http.Client, readInp ReadDomainInfoInput) (*ReadDomainInfoOutput, error) {

	readUrl := url.ReadFmcDomainInfo(readInp.FmcHost)

	req := client.NewGet(ctx, readUrl)

	var outp ReadDomainInfoOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
