package users_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth/role"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/users"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	netHttp "net/http"
	"strconv"
	"testing"
	"time"
)

// Function to generate users
func generateUsers(num int) []users.UserDetails {
	var createdUsers []users.UserDetails
	for i := 1; i <= num; i++ {
		uid := "uid" + strconv.Itoa(i)       // Generate unique UID
		username := "user" + strconv.Itoa(i) // Generate usernames like user1, user2, etc.
		roles := []string{"ROLE_USER"}       // Assign a default role; you can modify this as needed
		apiOnlyUser := i%2 == 0              // Example: alternate between true/false for ApiOnlyUser

		createdUsers = append(createdUsers, users.UserDetails{
			Uid:         uid,
			Username:    username,
			Roles:       roles,
			ApiOnlyUser: apiOnlyUser,
		})
	}
	return createdUsers
}

// the create test also tests read!
func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully create users in MSP-managed tenant", func(t *testing.T) {
		httpmock.Reset()
		var managedTenantUid = uuid.New().String()
		var createInp = users.MspUsersInput{
			TenantUid: managedTenantUid,
			Users: []users.UserDetails{
				{Username: "apples@bananas.com", Roles: []string{string(role.SuperAdmin)}, ApiOnlyUser: false},
				{Username: "api-only-user", Roles: []string{string(role.ReadOnly)}, ApiOnlyUser: true},
			},
		}

		var usersInCdoTenant = generateUsers(250)
		var usersWithIds []users.UserDetails
		for _, user := range createInp.Users {
			userWithId := users.UserDetails{
				Uid:         uuid.New().String(),
				Username:    user.Username,
				Roles:       user.Roles,
				ApiOnlyUser: user.ApiOnlyUser,
			}
			usersInCdoTenant = append(usersInCdoTenant, userWithId)
			usersWithIds = append(usersWithIds, userWithId)
		}
		firstUserPage := users.UserPage{Items: usersInCdoTenant[:200], Count: len(usersInCdoTenant), Limit: 200, Offset: 0}
		secondUserPage := users.UserPage{Items: usersInCdoTenant[200:], Count: len(usersInCdoTenant), Limit: 200, Offset: 200}
		var transactionUid = uuid.New().String()
		var inProgressTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USERS_TO_TENANT,
			Status:          transactionstatus.IN_PROGRESS,
		}
		var doneTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USERS_TO_TENANT,
			Status:          transactionstatus.DONE,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, inProgressTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			inProgressTransaction.PollingUrl,
			httpmock.NewJsonResponderOrPanic(200, doneTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users?limit=200&offset=0", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, firstUserPage),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users?limit=200&offset=200", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, secondUserPage),
		)

		actual, err := users.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), createInp)

		assert.NotNil(t, actual, "Created users should have not been nil")
		assert.Nil(t, err, "Created users operation should have not been an error")
		assert.Equal(t, usersWithIds, *actual, "Created users operation should have been the same as the created tenant")
	})

	t.Run("user creation transaction fails", func(t *testing.T) {
		httpmock.Reset()
		var managedTenantUid = uuid.New().String()
		var createInp = users.MspUsersInput{
			TenantUid: managedTenantUid,
			Users: []users.UserDetails{
				{Username: "apples@bananas.com", Roles: []string{string(role.SuperAdmin)}, ApiOnlyUser: false},
				{Username: "api-only-user", Roles: []string{string(role.ReadOnly)}, ApiOnlyUser: true},
			},
		}
		var transactionUid = uuid.New().String()
		var inProgressTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USERS_TO_TENANT,
			Status:          transactionstatus.IN_PROGRESS,
		}
		var errorTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USERS_TO_TENANT,
			Status:          transactionstatus.ERROR,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, inProgressTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			inProgressTransaction.PollingUrl,
			httpmock.NewJsonResponderOrPanic(200, errorTransaction),
		)

		actual, err := users.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), createInp)

		assert.Nil(t, actual, "Created users should be nil")
		assert.NotNil(t, err, "Created users in tenant operation should have an error")
		assert.Equal(t, users.CreateError{
			Err:               publicapi.NewTransactionErrorFromTransaction(errorTransaction),
			CreatedResourceId: &managedTenantUid,
		}, *err, "created transaction error does not match")
	})

	t.Run("user creation API call fails", func(t *testing.T) {
		httpmock.Reset()
		var managedTenantUid = uuid.New().String()
		var createInp = users.MspUsersInput{
			TenantUid: managedTenantUid,
			Users: []users.UserDetails{
				{Username: "apples@bananas.com", Roles: []string{string(role.SuperAdmin)}, ApiOnlyUser: false},
				{Username: "api-only-user", Roles: []string{string(role.ReadOnly)}, ApiOnlyUser: true},
			},
		}
		var transactionUid = uuid.New().String()
		var errorTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USERS_TO_TENANT,
			Status:          transactionstatus.ERROR,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, errorTransaction),
		)
		actual, err := users.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), createInp)

		assert.Nil(t, actual, "Created users in tenant should have not been nil")
		assert.NotNil(t, err, "Created users in tenant operation should have not been an error")
	})
}
