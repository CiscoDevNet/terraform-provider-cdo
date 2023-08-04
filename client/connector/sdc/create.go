package sdc

import (
	"context"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
	"github.com/cisco-lockhart/go-client/user"
)

type CreateInput struct {
	Name string
}

func NewCreateInput(sdcName string) *CreateInput {
	return &CreateInput{
		Name: sdcName,
	}
}

type createRequestBody struct {
	Name                string `json:"name"`
	OnPremLarConfigured bool   `json:"onPremLarConfigured"`
}

type createRequestOutput struct {
	Uid                      string `json:"uid"`
	Name                     string `json:"name"`
	Status                   string `json:"status"`
	State                    string `json:"state"`
	TenantUid                string `json:"tenantUid"`
	ServiceConnectivityState string `json:"serviceConnectivityState"`
}

type UpdateOutput struct {
	*createRequestOutput
	BootstrapData string
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*UpdateOutput, error) {

	client.Logger.Println("create SDC")

	// 1. create sdc device
	url := url.CreateSdc(client.BaseUrl())
	body := createRequestBody{
		Name:                createInp.Name,
		OnPremLarConfigured: true, // TODO: when will this be false? See related: https://github.com/cisco-lockhart/eos/blob/4d2a8e7414073ac466b47647e834feb60abdef79/client/app/sdc/sdc.controller.js#L177C1
	}
	req := client.NewPost(ctx, url, body)

	var createOutp createRequestOutput
	if err := req.Send(&createOutp); err != nil {
		return &UpdateOutput{}, err
	}

	// 2. generate bootstrap data
	// get user data from authentication service
	userToken, err := user.GetToken(ctx, client, user.NewGetTokenInput())
	if err != nil {
		return &UpdateOutput{}, err
	}
	host, err := client.Host()
	if err != nil {
		return &UpdateOutput{}, err
	}
	bootstrapData := computeBootstrapData(
		createInp.Name, userToken.AccessToken, userToken.TenantName, client.BaseUrl(), host,
	)

	// 3. done!
	return &UpdateOutput{
		createRequestOutput: &createOutp,
		BootstrapData:       bootstrapData,
	}, nil
}
