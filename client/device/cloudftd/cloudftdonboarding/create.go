package cloudftdonboarding

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type CreateInput struct {
	FtdUid string `json:"ftdUid"`
}

func NewCreateInput(ftdId string) CreateInput {
	return CreateInput{
		FtdUid: ftdId,
	}
}

type CreateOutput = cloudftd.ReadOutput

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating cloud ftd onboarding")

	createUrl := url.RegisterFtd(client.BaseUrl())

	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		createUrl,
		createInp,
	)
	if err != nil {
		return nil, err
	}
	transaction, err = publicapi.WaitForTransactionToFinishWithDefaults(
		ctx,
		client,
		transaction,
		"Waiting for FTD onboarding to finish...",
	)
	if err != nil {
		return nil, err
	}

	return cloudftd.ReadByUid(ctx, client, cloudftd.NewReadByUidInput(createInp.FtdUid))
}
