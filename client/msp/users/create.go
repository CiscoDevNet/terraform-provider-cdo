package users

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func Create(ctx context.Context, client http.Client, createInp MspCreateUsersInput) (*[]UserDetails, *CreateError) {
	client.Logger.Printf("Creating %d users in %s\n", len(createInp.Users), createInp.TenantUid)
	createUrl := url.CreateUsersInMspManagedTenant(client.BaseUrl(), createInp.TenantUid)
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
		fmt.Sprintf("Waiting for users to be created and added to MSP-managed tenant %s...", createInp.TenantUid),
	)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}

	return &createInp.Users, nil
}
