package sdc

import (
	"context"

	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateInput struct {
	Uid  string `json:"-"`
	Name string `json:"name"`
}

func NewUpdateInput(uid string, name string) UpdateInput {
	return UpdateInput{
		Uid:  uid,
		Name: name,
	}
}

type UpdateOutput struct {
	*updateRequestOutput
	BootstrapData string
}

type updateRequestOutput struct {
	Uid                      string `json:"uid"`
	Name                     string `json:"name"`
	Status                   string `json:"status"`
	State                    string `json:"state"`
	TenantUid                string `json:"tenantUid"`
	ServiceConnectivityState string `json:"serviceConnectivityState"`
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	url := url.UpdateSdc(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, url, updateInp)

	var updateOutp updateRequestOutput
	if err := req.Send(&updateOutp); err != nil {
		return &UpdateOutput{}, nil
	}

	bootstrapData, err := generateBootstrapData(ctx, client, updateOutp.Name)
	if err != nil {
		return &UpdateOutput{}, nil
	}

	// 3. done!
	return &UpdateOutput{
		updateRequestOutput: &updateOutp,
		BootstrapData:       bootstrapData,
	}, nil
}
