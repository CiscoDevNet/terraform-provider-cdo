package connector

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateInput struct {
	Uid  string `json:"-"`
	Name string `json:"name,omitempty"`
}

func NewUpdateInput(uid string, name string) UpdateInput {
	return UpdateInput{
		Uid:  uid,
		Name: name,
	}
}

type UpdateOutput struct {
	*UpdateRequestOutput
	BootstrapData string
}

type UpdateRequestOutput struct {
	Uid                      string `json:"uid"`
	Name                     string `json:"name"`
	Status                   string `json:"status"`
	State                    string `json:"state"`
	TenantUid                string `json:"tenantUid"`
	ServiceConnectivityState string `json:"serviceConnectivityState"`
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	url := url.UpdateConnector(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, url, updateInp)

	var updateOutp UpdateRequestOutput
	if err := req.Send(&updateOutp); err != nil {
		return &UpdateOutput{}, err
	}

	bootstrapData, err := generateBootstrapData(ctx, client, updateOutp.Name)
	if err != nil {
		return &UpdateOutput{}, err
	}

	// 3. done!
	return &UpdateOutput{
		UpdateRequestOutput: &updateOutp,
		BootstrapData:       bootstrapData,
	}, nil
}
