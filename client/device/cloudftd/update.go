package cloudftd

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateInput struct {
	Uid  string
	Name string
}

func NewUpdateInput(uid, name string) UpdateInput {
	return UpdateInput{
		Uid:  uid,
		Name: name,
	}
}

type updateRequestBody struct {
	Name string `json:"name"`
}

type UpdateOutput = ReadByUidOutput

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	updateUrl := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)
	updateBody := updateRequestBody{
		Name: updateInp.Name,
	}
	req := client.NewPut(ctx, updateUrl, updateBody)
	var updateOutp UpdateOutput
	if err := req.Send(&updateOutp); err != nil {
		return nil, err
	}

	return &updateOutp, nil

}
