package users_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/users"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	netHttp "net/http"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully delete users in MSP-managed tenant", func(t *testing.T) {
		httpmock.Reset()
		managedTenantUid := uuid.New().String()
		deleteInp := users.MspDeleteUsersInput{
			TenantUid: managedTenantUid,
			Usernames: []string{"user1@example.com", "api-only-user", "user3@example.com"},
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
			Type:            transactiontype.MSP_DELETE_USERS_FROM_TENANT,
			Status:          transactionstatus.DONE,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users/delete", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, inProgressTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			inProgressTransaction.PollingUrl,
			httpmock.NewJsonResponderOrPanic(200, doneTransaction),
		)

		actual, err := users.Delete(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), deleteInp)

		assert.Nil(t, actual, "Deletion output should be nil")
		assert.Nil(t, err, "Deletion error should be nil")
	})

	t.Run("transaction to delete users in MSP-managed tenant fails", func(t *testing.T) {
		httpmock.Reset()
		managedTenantUid := uuid.New().String()
		deleteInp := users.MspDeleteUsersInput{
			TenantUid: managedTenantUid,
			Usernames: []string{"user1@example.com", "api-only-user", "user3@example.com"},
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
			Type:            transactiontype.MSP_DELETE_USERS_FROM_TENANT,
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
			Type:            transactiontype.MSP_DELETE_USERS_FROM_TENANT,
			Status:          transactionstatus.ERROR,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users/delete", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, inProgressTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			inProgressTransaction.PollingUrl,
			httpmock.NewJsonResponderOrPanic(200, errorTransaction),
		)

		actual, err := users.Delete(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), deleteInp)

		assert.Nil(t, actual, "Deletion output should be nil")
		assert.NotNil(t, err, "Deletion error should be nil")
		assert.Equal(t, err.Error(), fmt.Sprintf("error: transaction failed, uid=%s, status=ERROR, message=, details=map[]", transactionUid))
	})
}
