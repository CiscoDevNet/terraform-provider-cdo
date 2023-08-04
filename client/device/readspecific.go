package device

import (
	"context"

	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/url"
)

type ReadSpecificInput struct {
	Uid string `json:"uid"`
}

type ReadSpecificOutput struct {
	SpecificUid string `json:"uid"`
	State       string `json:"state"`
	Namespace   string `json:"namespace"`
	Type        string `json:"type"`
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

func ReadSpecific(ctx context.Context, client http.Client, readReq ReadSpecificInput) (*ReadSpecificOutput, error) {

	client.Logger.Println("reading specific device")

	req := NewReadSpecificRequest(ctx, client, readReq)

	var outp ReadSpecificOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
