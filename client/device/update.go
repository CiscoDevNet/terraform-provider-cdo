package device

import (
	"context"

	"github.com/CiscoDevnet/go-client/internal/http"
	"github.com/CiscoDevnet/go-client/internal/url"
)

type UpdateInput struct {
	Uid string `json:"-"`

	Name             string `json:"name,omitempty"`
	IgnoreCertifcate bool   `json:"ignoreCertificate,omitempty"`
}

type UpdateOutput = ReadOutput

func NewUpdateInput(uid string, name string, ignoreCertificate bool) *UpdateInput {
	return &UpdateInput{
		Uid:              uid,
		Name:             name,
		IgnoreCertifcate: ignoreCertificate,
	}
}

func NewUpdateRequest(ctx context.Context, client http.Client, updateReq UpdateInput) *http.Request {
	url := url.UpdateDevice(client.BaseUrl(), updateReq.Uid)

	req := client.NewPut(ctx, url, updateReq)

	return req
}

func Update(ctx context.Context, client http.Client, updateReq UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating device")

	req := NewUpdateRequest(ctx, client, updateReq)

	var outp UpdateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
