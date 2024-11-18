package usergroups

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func Delete(ctx context.Context, client http.Client, tenantUid string, deleteInp *MspManagedUserGroupDeleteInput) (interface{}, error) {
	client.Logger.Printf("Deleting %d user groups in %s\n", len(deleteInp.UserGroupUids), tenantUid)
	deleteUrl := url.DeleteUserGroupsInMspManagedTenant(client.BaseUrl(), tenantUid)
	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		deleteUrl,
		deleteInp,
	)
	if err != nil {
		return nil, err
	}

	_, err = publicapi.WaitForTransactionToFinishWithDefaults(
		ctx,
		client,
		transaction,
		fmt.Sprintf("Waiting for user groups to be deleted from MSP-managed tenant %s...", tenantUid),
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
