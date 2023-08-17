package genericssh

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type UpdateInput struct {
	Uid  string `json:"-"`
	Name string `json:"name,omitempty"`
}

type UpdateOutput = device.ReadOutput

func NewUpdateInput(uid string, name string) *UpdateInput {
	return &UpdateInput{
		Uid:  uid,
		Name: name,
	}
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating generic ssh")

	updateUrl := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, updateUrl, updateInp)

	var outp UpdateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
