package duoadminpanel

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"time"
)

type CreateInput struct {
	Name           string   `json:"name"`
	Host           string   `json:"host"`
	IntegrationKey string   `json:"integrationKey"`
	SecretKey      string   `json:"secretKey"`
	Labels         []string `json:"labels"`
}

type CreateOutput = ReadOutput

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating duo admin panel")

	createUrl := url.CreateDuoAdminPanel(client.BaseUrl())

	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		createUrl,
		createInp,
	)
	if err != nil {
		return nil, err
	}
	transaction, err = publicapi.PollTransaction(
		ctx,
		client,
		transaction,
		retry.NewOptionsBuilder().
			Logger(client.Logger).
			Timeout(5*time.Minute).
			Retries(-1).
			EarlyExitOnError(true).
			Message("Waiting for Duo Admin Panel to onboard...").
			Delay(1*time.Second).
			Build(),
	)
	if err != nil {
		return nil, err
	}

	return ReadByUid(ctx, client, NewReadByUidInput(transaction.EntityUid))
}
