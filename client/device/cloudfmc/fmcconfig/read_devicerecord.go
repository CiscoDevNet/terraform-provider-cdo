package fmcconfig

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/fmcconfig"
)

type ReadDeviceRecordInput struct {
	FmcDomainUid    string
	FmcHostname     string
	DeviceRecordUid string
}

func NewReadDeviceRecordInput(fmcDomainUid, fmcHostname, deviceRecordUid string) ReadDeviceRecordInput {
	return ReadDeviceRecordInput{
		FmcDomainUid:    fmcDomainUid,
		FmcHostname:     fmcHostname,
		DeviceRecordUid: deviceRecordUid,
	}
}

type ReadDeviceRecordOutput = fmcconfig.DeviceRecord

func ReadDeviceRecord(ctx context.Context, client http.Client, readInp ReadDeviceRecordInput) (*ReadDeviceRecordOutput, error) {

	readUrl := url.ReadFmcDeviceRecord(client.BaseUrl(), readInp.FmcDomainUid, readInp.DeviceRecordUid)

	req := client.NewGet(ctx, readUrl)
	req.Header.Add("Fmc-Hostname", readInp.FmcHostname)

	var readOutp fmcconfig.DeviceRecord
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}
