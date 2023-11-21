package fmcplatform

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcdomain"
)

type ReadDomainInfoInput struct {
	FmcHost string
}

func NewReadDomainInfoInput(fmcHost string) ReadDomainInfoInput {
	return ReadDomainInfoInput{
		FmcHost: fmcHost,
	}
}

type ReadDomainInfoOutput = fmcdomain.Info

func ReadFmcDomainInfo(ctx context.Context, client http.Client, readInp ReadDomainInfoInput) (*ReadDomainInfoOutput, error) {

	client.Logger.Println("reading FMC domain info")

	readUrl := url.ReadFmcDomainInfo(readInp.FmcHost)

	req := client.NewGet(ctx, readUrl)

	var outp ReadDomainInfoOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
