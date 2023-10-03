package fmcconfig

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcconfig"
)

type ReadAllDeviceRecordsInput struct {
	FmcDomainUid string
	FmcHostname  string
}

func NewReadAllDeviceRecordsInput(fmcDomainUid string, fmcHostname string) ReadAllDeviceRecordsInput {
	return ReadAllDeviceRecordsInput{
		FmcDomainUid: fmcDomainUid,
		FmcHostname:  fmcHostname,
	}
}

type ReadAllDeviceRecordsOutput = fmcconfig.AllDeviceRecords

func ReadAllDeviceRecords(ctx context.Context, client http.Client, readInp ReadAllDeviceRecordsInput) (*ReadAllDeviceRecordsOutput, error) {

	readUrl := url.ReadFmcAllDeviceRecords(client.BaseUrl(), readInp.FmcDomainUid)

	req := client.NewGet(ctx, readUrl)
	req.Header.Add("Fmc-Hostname", readInp.FmcHostname)

	var readOutp ReadAllDeviceRecordsOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}
