package cloudftd

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
)

type UpdateInput struct {
	Uid  string
	Name string
	Tags tags.Type
}

func NewUpdateInput(uid, name string, tags tags.Type) UpdateInput {
	return UpdateInput{
		Uid:  uid,
		Name: name,
		Tags: tags,
	}
}

type updateRequestBody struct {
	Name string    `json:"name"`
	Tags tags.Type `json:"tags"`
}

type UpdateOutput = ReadOutput

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	updateUrl := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)
	updateBody := updateRequestBody{
		Name: updateInp.Name,
		Tags: updateInp.Tags,
	}
	req := client.NewPut(ctx, updateUrl, updateBody)
	var updateOutp UpdateOutput
	if err := req.Send(&updateOutp); err != nil {
		return nil, err
	}

	return &updateOutp, nil

}
