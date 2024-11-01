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

	_, err = publicapi.WaitForTransactionToFinishWithDefaults(
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

func RevokeApiToken(ctx context.Context, client http.Client, revokeInput MspRevokeApiTokenInput) (interface{}, error) {
	revokeTokenUrl := url.RevokeApiTokenUsingPublicApi(client.BaseUrl())
	client.Logger.Printf("Revoking api token at %s\n", revokeTokenUrl)
	req := client.NewPost(ctx, revokeTokenUrl, nil)
	// overwrite token in header with API token for the user that we are revoking
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", revokeInput.ApiToken))
	if err := req.Send(&struct{}{}); err != nil {
		return nil, err
	}

	return nil, nil
}
