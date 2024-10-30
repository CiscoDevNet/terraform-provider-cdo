package users

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func Delete(ctx context.Context, client http.Client, deleteInp MspDeleteUsersInput) (interface{}, error) {
	client.Logger.Printf("Deleting %d users in %s\n", len(deleteInp.Usernames), deleteInp.TenantUid)
	deleteUrl := url.DeleteUsersInMspManagedTenant(client.BaseUrl(), deleteInp.TenantUid)
	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		deleteUrl,
		deleteInp,
	)
	if err != nil {
		return nil, err
	}

	transaction, err = publicapi.WaitForTransactionToFinishWithDefaults(
		ctx,
		client,
		transaction,
		fmt.Sprintf("Waiting for users to be deleted from MSP-managed tenant %s...", deleteInp.TenantUid),
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
