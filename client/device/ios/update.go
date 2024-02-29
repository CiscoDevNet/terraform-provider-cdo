package ios

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateInput struct {
	Uid  string    `json:"-"`
	Name string    `json:"name"`
	Tags tags.Type `json:"tags"`
}

type UpdateOutput = device.UpdateOutput

func NewUpdateInput(uid string, name string, tags tags.Type) *UpdateInput {
	return &UpdateInput{
		Uid:  uid,
		Name: name,
		Tags: tags,
	}
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating ios device")

	url := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, url, updateInp)

	var outp UpdateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
