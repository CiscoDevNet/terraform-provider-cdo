package tenants

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func Create(ctx context.Context, client http.Client, createInp MspCreateTenantInput) (*MspTenantOutput, *CreateError) {
	client.Logger.Println("Creating tenant for CDO")
	createUrl := url.CreateMspManagedTenant(client.BaseUrl())

	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		createUrl,
		createInp,
	)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}
	transaction, err = publicapi.WaitForTransactionToFinishWithDefaults(
		ctx,
		client,
		transaction,
		"Waiting for tenant to be created and added to MSP portal...",
	)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}

	readOut, err := Read(ctx, client, ReadByUidInput{Uid: transaction.EntityUid})
	client.Logger.Println("Created tenant for CDO")
	if err == nil {
		return readOut, nil
	} else {
		return readOut, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}

}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}
