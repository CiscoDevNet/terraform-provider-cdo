package sdc

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type ReadByUidInput struct {
	SdcUid string
}

type ReadByNameInput struct {
	SdcName string
}

type ReadOutput struct {
	Uid        string          `json:"uid"`
	Name       string          `json:"name"`
	DefaultSdc bool            `json:"defaultLar"`
	Cdg        bool            `json:"cdg"`
	TenantUid  string          `json:"tenantUid"`
	PublicKey  model.PublicKey `json:"larPublicKey"`
}

func NewReadByUidInput(sdcUid string) *ReadByUidInput {
	return &ReadByUidInput{
		SdcUid: sdcUid,
	}
}

func NewReadByNameInput(sdcName string) *ReadByNameInput {
	return &ReadByNameInput{
		SdcName: sdcName,
	}
}

func newReadByUidRequest(ctx context.Context, client http.Client, readInp ReadByUidInput) *http.Request {

	url := url.ReadSdcByUid(client.BaseUrl(), readInp.SdcUid)

	req := client.NewGet(ctx, url)

	return req
}

func newReadByNameRequest(ctx context.Context, client http.Client, readInp ReadByNameInput) *http.Request {

	url := url.ReadSdcByName(client.BaseUrl(), readInp.SdcName)

	req := client.NewGet(ctx, url)

	return req
}

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadOutput, error) {

	client.Logger.Println("reading sdc")

	req := newReadByUidRequest(ctx, client, readInp)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}

func ReadByName(ctx context.Context, client http.Client, readInp ReadByNameInput) (*ReadOutput, error) {

	client.Logger.Println("reading sdc by name")

	req := newReadByNameRequest(ctx, client, readInp)

	var arrayOutp []ReadOutput
	if err := req.Send(&arrayOutp); err != nil {
		return nil, err
	}

	if len(arrayOutp) == 0 {
		return nil, fmt.Errorf("no SDC found")
	}

	if len(arrayOutp) > 1 {
		return nil, fmt.Errorf("multiple SDCs found with the name: %s", readInp.SdcName)
	}

	outp := arrayOutp[0]
	return &outp, nil
}
