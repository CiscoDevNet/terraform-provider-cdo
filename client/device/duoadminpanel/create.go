package duoadminpanel

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/publicapilabels"
)

type CreateInput struct {
	Name           string               `json:"name"`
	Host           string               `json:"host"`
	IntegrationKey string               `json:"integrationKey"`
	SecretKey      string               `json:"secretKey"`
	Labels         publicapilabels.Type `json:"labels"`
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
		_, _ = Delete(ctx, client, DeleteInput{Uid: transaction.TransactionUid})
		return nil, err
	}
	transaction, err = publicapi.WaitForTransactionToFinishWithDefaults(
		ctx,
		client,
		transaction,
		"Waiting for Duo Admin Panel to onboard...",
	)
	if err != nil {
		_, _ = Delete(ctx, client, DeleteInput{Uid: transaction.TransactionUid})
		return nil, err
	}

	return ReadByUid(ctx, client, NewReadByUidInput(transaction.EntityUid))
}
