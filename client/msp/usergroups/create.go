package usergroups

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func Create(ctx context.Context, client http.Client, tenantUid string, userGroupsInput *[]MspManagedUserGroupInput) (*[]MspManagedUserGroup, *CreateError) {
	client.Logger.Printf("Creating %d user groups in %s\n", len(*userGroupsInput), tenantUid)
	createUrl := url.CreateUserGroupsInMspManagedTenant(client.BaseUrl(), tenantUid)
	transaction, err := publicapi.TriggerTransaction(ctx, client, createUrl, userGroupsInput)
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
		fmt.Sprintf("Waiting for users to be created and added to MSP-managed tenant %s...", tenantUid),
	)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}

	readUserGroupDetails, err := ReadCreatedUserGroupsInTenant(ctx, client, tenantUid, userGroupsInput)
	if err != nil {
		client.Logger.Println("Failed to read users from tenant after creation")
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}
	return readUserGroupDetails, nil
}
