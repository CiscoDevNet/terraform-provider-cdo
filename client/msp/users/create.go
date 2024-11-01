package users

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

func Create(ctx context.Context, client http.Client, createInp MspUsersInput) (*[]UserDetails, *CreateError) {
	client.Logger.Printf("Creating %d users in %s\n", len(createInp.Users), createInp.TenantUid)
	createUrl := url.CreateUsersInMspManagedTenant(client.BaseUrl(), createInp.TenantUid)
	var userDetailsPublicApiInput []UserDetailsPublicApiInput
	for _, user := range createInp.Users {
		userDetailsPublicApiInput = append(userDetailsPublicApiInput, UserDetailsPublicApiInput{
			Username:    user.Username,
			Role:        user.Roles[0],
			ApiOnlyUser: user.ApiOnlyUser,
		})
	}
	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		createUrl,
		MspUsersPublicApiInput{
			TenantUid: createInp.TenantUid,
			Users:     userDetailsPublicApiInput,
		},
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

	readUserDetrails, err := ReadCreatedUsersInTenant(ctx, client, createInp)
	if err != nil {
		client.Logger.Println("Failed to read users from tenant after creation")
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}
	return readUserDetrails, nil
}

func GenerateApiToken(ctx context.Context, client http.Client, generateApiTokenInp MspGenerateApiTokenInput) (*MspGenerateApiTokenOutput, error) {
	client.Logger.Printf("Generating API token for user %s in tenant %s", generateApiTokenInp.UserUid, generateApiTokenInp.TenantUid)
	genApiTokenUrl := url.GenerateApiTokenForUserInMspManagedTenant(client.BaseUrl(), generateApiTokenInp.TenantUid, generateApiTokenInp.UserUid)
	var mspApiTokenOutput MspGenerateApiTokenOutput
	req := client.NewPost(ctx, genApiTokenUrl, nil)
	if err := req.Send(&mspApiTokenOutput); err != nil {
		return nil, err
	}
	return &mspApiTokenOutput, nil
}
