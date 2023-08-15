package iosconfig

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type ReadInput struct {
	SpecificUid string
}

type ReadOutput struct {
	Uid   string `json:"uid"`
	State string `json:"state"`
}

func NewReadInput(specificUid string) *ReadInput {
	return &ReadInput{
		SpecificUid: specificUid,
	}
}

func NewReadRequest(ctx context.Context, client http.Client, readReq ReadInput) *http.Request {

	readUrl := url.ReadDevice(client.BaseUrl(), readReq.SpecificUid)

	req := client.NewGet(ctx, readUrl)

	return req
}

func Read(ctx context.Context, client http.Client, readReq ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading iosconfig")

	req := NewReadRequest(ctx, client, readReq)

	var outp ReadOutput
	err := req.Send(&outp)
	if err != nil {
		return nil, err
	}

	return &outp, nil
}
