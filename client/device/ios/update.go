package ios

import (
	"context"

	"github.com/cisco-lockhart/go-client/device"
	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type UpdateInput struct {
	Uid  string `json:"-"`
	Name string `json:"name"`
}

type UpdateOutput = device.UpdateOutput

func NewUpdateInput(uid string, name string) *UpdateInput {
	return &UpdateInput{
		Uid:  uid,
		Name: name,
	}
}

func NewUpdateRequest(ctx context.Context, client http.Client, updateInp UpdateInput) *http.Request {

	url := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, url, updateInp)

	return req
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating ios device")

	req := NewUpdateRequest(ctx, client, updateInp)

	var outp UpdateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
