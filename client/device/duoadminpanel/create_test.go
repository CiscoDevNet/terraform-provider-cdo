package duoadminpanel_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestDuoAdminPanelCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	createInput := duoadminpanel.CreateInput{
		Name:           "test-name",
		Host:           "test-host",
		IntegrationKey: "test-int-key",
		SecretKey:      "test-secret-key",
		Labels:         []string{"lab1", "lab2", "lab3"},
	}

	doneTransaction := transaction.Type{
		TransactionUid:  "test-TransactionUid",
		TenantUid:       "test-TenantUid",
		EntityUid:       "test-EntityUid",
		EntityUrl:       "test-EntityUrl",
		PollingUrl:      "test-PollingUrl",
		SubmissionTime:  "test-SubmissionTime",
		LastUpdatedTime: "test-LastUpdatedTime",
		Type:            transactiontype.ONBOARD_DUO_ADMIN_PANEL,
		Status:          transactionstatus.DONE,
	}

	errorTransaction := transaction.Type{
		TransactionUid:  "test-TransactionUid",
		TenantUid:       "test-TenantUid",
		EntityUid:       "test-EntityUid",
		EntityUrl:       "test-EntityUrl",
		PollingUrl:      "test-PollingUrl",
		SubmissionTime:  "test-SubmissionTime",
		LastUpdatedTime: "test-LastUpdatedTime",
		Type:            transactiontype.ONBOARD_DUO_ADMIN_PANEL,
		Status:          transactionstatus.ERROR,
		ErrorMessage:    "test-ErrorMessage",
		ErrorDetails: map[string]string{
			"test-key": "test-details",
		},
	}

	createdDevice := duoadminpanel.ReadOutput{
		Uid:  doneTransaction.EntityUid,
		Name: createInput.Name,
		Tags: tags.New(createInput.Labels...),
	}

	expectedCreateOutput := duoadminpanel.CreateOutput{
		Uid:   createdDevice.Uid,
		Name:  createdDevice.Name,
		State: createdDevice.State,
		Tags:  createdDevice.Tags,
	}

	baseUrl := "https://test.cisco.com"

	testCases := []struct {
		testName   string
		input      duoadminpanel.CreateInput
		setupFunc  func(input duoadminpanel.CreateInput)
		assertFunc func(output *duoadminpanel.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully onboards Duo Admin Panel",
			input:    createInput,

			setupFunc: func(input duoadminpanel.CreateInput) {
				configureTransactionToReturn(url.CreateDuoAdminPanel(baseUrl), doneTransaction)
				configureReadDeviceToReturn(url.ReadDevice(baseUrl, createdDevice.Uid), createdDevice)
			},

			assertFunc: func(actualOutput *duoadminpanel.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, actualOutput)
				assert.Equal(t, expectedCreateOutput, *actualOutput)
			},
		},
		{
			testName: "fails onboards Duo Admin Panel if transaction fails",
			input:    createInput,

			setupFunc: func(input duoadminpanel.CreateInput) {
				configureTransactionToReturn(url.CreateDuoAdminPanel(baseUrl), errorTransaction)
				configureReadDeviceToReturn(url.ReadDevice(baseUrl, createdDevice.Uid), createdDevice)
			},

			assertFunc: func(actualOutput *duoadminpanel.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, actualOutput)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, errorTransaction.ErrorMessage)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := duoadminpanel.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func configureTransactionToReturn(url string, t transaction.Type) {
	httpmock.RegisterResponder(http.MethodPost, url, httpmock.NewJsonResponderOrPanic(202, t))
}

func configureReadDeviceToReturn(url string, device duoadminpanel.ReadOutput) {
	httpmock.RegisterResponder(http.MethodGet, url, httpmock.NewJsonResponderOrPanic(202, device))
}
