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
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully create users in MSP-managed tenant", func(t *testing.T) {
		httpmock.Reset()
		var managedTenantUid = uuid.New().String()
		var createInp = users.MspCreateUsersInput{
			TenantUid: managedTenantUid,
			Users: []users.UserInput{
				{Username: "apples@bananas.com", Role: string(role.SuperAdmin), ApiOnlyUser: false},
				{Username: "api-only-user", Role: string(role.ReadOnly), ApiOnlyUser: true},
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

		actual, err := users.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), createInp)

		assert.NotNil(t, actual, "Created users should have not been nil")
		assert.Nil(t, err, "Created users operation should have not been an error")
		assert.Equal(t, createInp.Users, *actual, "Created users operation should have been the same as the created tenant")
	})

	t.Run("user creation transaction fails", func(t *testing.T) {
		httpmock.Reset()
		var managedTenantUid = uuid.New().String()
		var createInp = users.MspCreateUsersInput{
			TenantUid: managedTenantUid,
			Users: []users.UserInput{
				{Username: "apples@bananas.com", Role: string(role.SuperAdmin), ApiOnlyUser: false},
				{Username: "api-only-user", Role: string(role.ReadOnly), ApiOnlyUser: true},
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
		var createInp = users.MspCreateUsersInput{
			TenantUid: managedTenantUid,
			Users: []users.UserInput{
				{Username: "apples@bananas.com", Role: string(role.SuperAdmin), ApiOnlyUser: false},
				{Username: "api-only-user", Role: string(role.ReadOnly), ApiOnlyUser: true},
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
