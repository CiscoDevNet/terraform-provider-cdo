package duoadminpanel_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestDuoAdminPanelCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testModel := internalTesting.NewRandomModel()

	createInput := testModel.DuoAdminPanelCreateInput()
	readOutput := testModel.DuoAdminPanelReadOutput()
	doneTransaction := testModel.CreateDoneTransaction(readOutput.Uid, transactiontype.ONBOARD_DUO_ADMIN_PANEL)
	errorTransaction := testModel.CreateErrorTransaction(readOutput.Uid, transactiontype.ONBOARD_DUO_ADMIN_PANEL)

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
				internalTesting.MockPostAccepted(url.CreateDuoAdminPanel(testModel.BaseUrl), doneTransaction)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
			},

			assertFunc: func(actualOutput *duoadminpanel.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, actualOutput)
				assert.Equal(t, readOutput, *actualOutput)
			},
		},
		{
			testName: "fails onboards Duo Admin Panel if transaction fails",
			input:    createInput,

			setupFunc: func(input duoadminpanel.CreateInput) {
				internalTesting.MockPostAccepted(url.CreateDuoAdminPanel(testModel.BaseUrl), errorTransaction)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
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
				*internalHttp.MustNewWithConfig(testModel.BaseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
