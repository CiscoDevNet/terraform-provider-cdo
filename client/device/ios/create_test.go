package ios_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestIosCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testModel := internalTesting.NewRandomModel()

	createInput := testModel.CreateIosInput()
	readOutput := testModel.ReadIosOutput()
	doneTransaction := testModel.CreateDoneTransaction(readOutput.Uid, transactiontype.ONBOARD_IOS)
	errorTransaction := testModel.CreateErrorTransaction(readOutput.Uid, transactiontype.ONBOARD_IOS)

	testCases := []struct {
		testName   string
		input      ios.CreateInput
		setupFunc  func(input ios.CreateInput)
		assertFunc func(output *ios.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully onboards IOS",
			input:    createInput,

			setupFunc: func(input ios.CreateInput) {
				internalTesting.MockPostAccepted(url.CreateIos(testModel.BaseUrl), doneTransaction)
				internalTesting.MockGetOk(url.ReadConnectorByUid(testModel.BaseUrl, testModel.CdgUid.String()), readOutput)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
			},

			assertFunc: func(actualOutput *ios.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, actualOutput)
				assert.Equal(t, readOutput, *actualOutput)
			},
		},
		{
			testName: "fails onboards Duo Admin Panel if transaction fails",
			input:    createInput,

			setupFunc: func(input ios.CreateInput) {
				internalTesting.MockPostError(url.CreateIos(testModel.BaseUrl), errorTransaction)
				internalTesting.MockGetOk(url.ReadConnectorByUid(testModel.BaseUrl, testModel.CdgUid.String()), readOutput)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
			},

			assertFunc: func(actualOutput *ios.CreateOutput, err error, t *testing.T) {
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

			output, err := ios.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(testModel.BaseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
