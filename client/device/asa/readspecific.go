package asa

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type ReadSpecificInput struct {
	Uid string `json:"uid"`
}

type ReadSpecificOutput struct {
	SpecificUid string                 `json:"uid"`
	State       state.Type             `json:"state"`
	Namespace   string                 `json:"namespace"`
	Type        string                 `json:"type"`
	Metadata    SpecificDeviceMetadata `json:"metadata"`
}

type SpecificDeviceMetadata struct {
	AsdmVersion string `json:"deviceManager"`
}

func NewReadSpecificInput(uid string) *ReadSpecificInput {
	return &ReadSpecificInput{
		Uid: uid,
	}
}

func NewReadSpecificRequest(ctx context.Context, client http.Client, readInp ReadSpecificInput) *http.Request {

	url := url.ReadSpecificDevice(client.BaseUrl(), readInp.Uid)

	req := client.NewGet(ctx, url)

	return req
}

func ReadSpecific(ctx context.Context, client http.Client, readInp ReadSpecificInput) (*ReadSpecificOutput, error) {

	client.Logger.Println("reading asa specific device")

	req := NewReadSpecificRequest(ctx, client, readInp)

	var outp ReadSpecificOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
