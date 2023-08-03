package sdc

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
	"github.com/cisco-lockhart/go-client/user"
)

type CreateInput struct {
	Name string
}

func NewCreateInput(larUid string) *CreateInput {
	return &CreateInput{
		Name: larUid,
	}
}

type createRequestBody struct {
	Name                string `json:"name"`
	OnPremLarConfigured bool   `json:"onPremLarConfigured"`
}

type createRequestOutput struct {
	Uid                      string `json:"uid"`
	Name                     bool   `json:"name"`
	Status                   string `json:"status"`
	State                    string `json:"state"`
	TenantUid                string `json:"tenantUid"`
	ServiceConnectivityState string `json:"serviceConnectivityState"`
}

type CreateOutput struct {
	*createRequestOutput
	BootstrapData string
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (CreateOutput, error) {

	client.Logger.Println("create SDC")

	// 1. create sdc device
	url := url.CreateSdc(client.BaseUrl())
	body := createRequestBody{
		Name:                createInp.Name,
		OnPremLarConfigured: true, // TODO: when will this be false? See also: https://github.com/cisco-lockhart/eos/blob/4d2a8e7414073ac466b47647e834feb60abdef79/client/app/sdc/sdc.controller.js#L177C1
	}
	req := client.NewPost(ctx, url, body)

	var createOutp createRequestOutput
	if err := req.Send(&createOutp); err != nil {
		return CreateOutput{}, err
	}

	// 2. generate bootstrap data
	// get user data from authentication service
	userToken, err := user.GetToken(ctx, client, user.NewGetTokenInput())
	if err != nil {
		return CreateOutput{}, err
	}
	host, err := client.Host()
	if err != nil {
		return CreateOutput{}, err
	}
	bootstrapData := computeBootstrapData(
		createInp.Name, userToken.AccessToken, userToken.TenantName, client.BaseUrl(), host,
	)

	// 3. done!
	return CreateOutput{
		createRequestOutput: &createOutp,
		BootstrapData:       bootstrapData,
	}, nil
}

func computeBootstrapData(sdcName, accessToken, tenantName, baseUrl, host string) string {
	bootstrapUrl := fmt.Sprintf("%s/sdc/bootstrap/%s/%s", baseUrl, tenantName, sdcName)

	rawBootstrapData := fmt.Sprintf("CDO_TOKEN=%q\nCDO_DOMAIN=%q\nCDO_TENANT=%q\nCDO_BOOTSTRAP_URL=%q\n", accessToken, host, tenantName, bootstrapUrl)

	bootstrapData := base64.StdEncoding.EncodeToString([]byte(rawBootstrapData))

	return bootstrapData
}
