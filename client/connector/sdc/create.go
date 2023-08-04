package sdc

import (
	"context"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
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

type CreateRequestOutput struct {
	Uid                      string `json:"uid"`
	Name                     string `json:"name"`
	Status                   string `json:"status"`
	State                    string `json:"state"`
	TenantUid                string `json:"tenantUid"`
	ServiceConnectivityState string `json:"serviceConnectivityState"`
}

type CreateOutput struct {
	*CreateRequestOutput
	BootstrapData string
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("create SDC")

	// 1. create sdc device
	url := url.CreateSdc(client.BaseUrl())
	body := createRequestBody{
		Name:                createInp.Name,
		OnPremLarConfigured: true, // TODO: when will this be false? See related: https://github.com/cisco-lockhart/eos/blob/4d2a8e7414073ac466b47647e834feb60abdef79/client/app/sdc/sdc.controller.js#L177C1
	}
	req := client.NewPost(ctx, url, body)

	var createOutp CreateRequestOutput
	if err := req.Send(&createOutp); err != nil {
		return &CreateOutput{}, err
	}

	// 2. generate bootstrap data
	// get user data from authentication service
	bootstrapData, err := generateBootstrapData(ctx, client, createInp.Name)
	if err != nil {
		return &CreateOutput{}, nil
	}

	// 3. done!
	return &CreateOutput{
		CreateRequestOutput: &createOutp,
		BootstrapData:       bootstrapData,
	}, nil
}
