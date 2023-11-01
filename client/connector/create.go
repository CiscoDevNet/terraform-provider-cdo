package connector

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"time"
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
	url := url.CreateConnector(client.BaseUrl())
	body := createRequestBody{
		Name:                createInp.Name,
		OnPremLarConfigured: true, // this is always true, because we no longer support cloud sdc
	}
	req := client.NewPost(ctx, url, body)

	var createOutp CreateRequestOutput
	if err := req.Send(&createOutp); err != nil {
		return &CreateOutput{}, err
	}

	// 1.5 poll until SDC has SQS/SNS setup properly to communicate with backend
	err := retry.Do(
		ctx,
		untilCommunicationQueueReady(ctx, client, *NewReadByUidInput(createOutp.Uid)),
		retry.NewOptionsBuilder().
			Message("waiting for SDC to be available").
			Retries(10).
			Logger(client.Logger).
			Delay(2*time.Second).
			EarlyExitOnError(true).
			Timeout(time.Minute). // typically a few seconds should be enough
			Build(),
	)
	if err != nil {
		return &CreateOutput{}, err
	}

	// 2. generate bootstrap data
	// get user data from authentication service
	bootstrapData, err := generateBootstrapData(ctx, client, createInp.Name)
	if err != nil {
		return &CreateOutput{}, err
	}

	// 3. done!
	return &CreateOutput{
		CreateRequestOutput: &createOutp,
		BootstrapData:       bootstrapData,
	}, nil
}

func untilCommunicationQueueReady(ctx context.Context, client http.Client, input ReadByUidInput) retry.Func {
	return func() (bool, error) {
		readOutp, err := ReadByUid(ctx, client, input)
		if err != nil {
			return false, err
		}
		return readOutp.IsCommunicationQueueReady, nil
	}
}
