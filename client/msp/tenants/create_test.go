package tenants_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/tenants"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	netHttp "net/http"
	"testing"
	"time"
)

const (
	baseUrl = "https://unittest.cdo.cisco.com"
)

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully create tenant", func(t *testing.T) {
		httpmock.Reset()
		var createInp = tenants.MspCreateTenantInput{
			Name: "test-tenant-name",
		}
		var entityUid = uuid.New().String()
		var transactionUid = uuid.New().String()
		var inProgressTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       entityUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + entityUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_CREATE_TENANT,
			Status:          transactionstatus.IN_PROGRESS,
		}
		var doneTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       entityUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + entityUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_CREATE_TENANT,
			Status:          transactionstatus.DONE,
		}
		var creationOutput = tenants.MspTenantOutput{
			Uid:         entityUid,
			Name:        createInp.Name,
			DisplayName: createInp.Name,
			Region:      "STAGING",
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/api/rest/v1/msp/tenants/create",
			httpmock.NewJsonResponderOrPanic(200, inProgressTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			inProgressTransaction.PollingUrl,
			httpmock.NewJsonResponderOrPanic(200, doneTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/api/rest/v1/msp/tenants/"+entityUid,
			httpmock.NewJsonResponderOrPanic(200, creationOutput),
		)

		actual, err := tenants.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), createInp)

		assert.NotNil(t, actual, "Created tenant should have not been nil")
		assert.Nil(t, err, "Created tenant operation should have not been an error")
		assert.Equal(t, creationOutput, *actual, "Created tenant operation should have been the same as the created tenant")
	})

	t.Run("tenant creation transaction fails", func(t *testing.T) {
		httpmock.Reset()
		var createInp = tenants.MspCreateTenantInput{
			Name: "test-tenant-name",
		}
		var entityUid = uuid.New().String()
		var transactionUid = uuid.New().String()
		var inProgressTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       entityUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + entityUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_CREATE_TENANT,
			Status:          transactionstatus.IN_PROGRESS,
		}
		var errorTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       entityUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + entityUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_CREATE_TENANT,
			Status:          transactionstatus.ERROR,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/api/rest/v1/msp/tenants/create",
			httpmock.NewJsonResponderOrPanic(200, inProgressTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			inProgressTransaction.PollingUrl,
			httpmock.NewJsonResponderOrPanic(200, errorTransaction),
		)

		actual, err := tenants.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), createInp)

		assert.Nil(t, actual, "Created tenant should have not been nil")
		assert.NotNil(t, err, "Created tenant operation should have not been an error")
		assert.Equal(t, tenants.CreateError{
			Err:               publicapi.NewTransactionErrorFromTransaction(errorTransaction),
			CreatedResourceId: &entityUid,
		}, *err, "created transaction error does not match")
	})

	t.Run("tenant creation API call fails", func(t *testing.T) {
		httpmock.Reset()
		var createInp = tenants.MspCreateTenantInput{
			Name: "test-tenant-name",
		}
		var entityUid = uuid.New().String()
		var transactionUid = uuid.New().String()
		var errorTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       entityUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + entityUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_CREATE_TENANT,
			Status:          transactionstatus.ERROR,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/api/rest/v1/msp/tenants/create",
			httpmock.NewJsonResponderOrPanic(200, errorTransaction),
		)
		actual, err := tenants.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), createInp)

		assert.Nil(t, actual, "Created tenant should have not been nil")
		assert.NotNil(t, err, "Created tenant operation should have not been an error")
	})
}
