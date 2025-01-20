package asa_test

import (
	"context"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/jarcoal/httpmock"
)

func TestAsaCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testModel := internalTesting.NewRandomModel()

	cdgReadOutput := testModel.CdgReadOutput()
	createInput := testModel.AsaCreateInput()
	readOutput := testModel.AsaReadOutput()
	readSpecificOutput := testModel.AsaReadSpecificDeviceOutput()
	createInputWithIncorrectSoftwareVersion := createInput
	createInputWithIncorrectSoftwareVersion.SoftwareVersion = readOutput.SoftwareVersion + "_different"
	createInputWithIncorrectAsdmVersion := createInput
	createInputWithIncorrectAsdmVersion.AsdmVersion = readOutput.AsdmVersion + "_different"
	doneTransaction := testModel.CreateDoneTransaction(readOutput.Uid, transactiontype.ONBOARD_ASA)
	errorTransaction := testModel.CreateErrorTransaction(readOutput.Uid, transactiontype.ONBOARD_ASA)

	testCases := []struct {
		testName   string
		input      asa.CreateInput
		setupFunc  func(input asa.CreateInput)
		assertFunc func(output *asa.ReadOutput, specificDeviceOutput *asa.ReadSpecificOutput, err *asa.CreateError, t *testing.T)
	}{
		{
			testName: "successfully onboards ASA",
			input:    createInput,

			setupFunc: func(input asa.CreateInput) {
				internalTesting.MockGetOk(url.ReadConnectorByUid(testModel.BaseUrl, cdgReadOutput.Uid), cdgReadOutput)
				internalTesting.MockPostAccepted(url.CreateAsa(testModel.BaseUrl), doneTransaction)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
				internalTesting.MockGetOk(url.ReadSpecificDevice(testModel.BaseUrl, readOutput.Uid), readSpecificOutput)
			},

			assertFunc: func(actualOutput *asa.ReadOutput, actualSpecificDeviceOutput *asa.ReadSpecificOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, actualOutput)
				assert.NotNil(t, actualSpecificDeviceOutput)
				assert.Equal(t, readOutput, *actualOutput)
			},
		},
		{
			testName: "fails to onboard ASA due to specified software version mismatch",
			input:    createInputWithIncorrectSoftwareVersion,
			setupFunc: func(input asa.CreateInput) {
				internalTesting.MockGetOk(url.ReadConnectorByUid(testModel.BaseUrl, cdgReadOutput.Uid), cdgReadOutput)
				internalTesting.MockPostAccepted(url.CreateAsa(testModel.BaseUrl), doneTransaction)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
				internalTesting.MockGetOk(url.ReadSpecificDevice(testModel.BaseUrl, readOutput.Uid), readSpecificOutput)
			},
			assertFunc: func(actualOutput *asa.ReadOutput, actualSpecificDeviceOutput *asa.ReadSpecificOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, actualOutput)
				assert.Nil(t, actualSpecificDeviceOutput)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, "ASA Software version mismatch.")
			},
		},
		{
			testName: "fails to onboard ASA due to specified ASDM version mismatch",
			input:    createInputWithIncorrectAsdmVersion,
			setupFunc: func(input asa.CreateInput) {
				internalTesting.MockGetOk(url.ReadConnectorByUid(testModel.BaseUrl, cdgReadOutput.Uid), cdgReadOutput)
				internalTesting.MockPostAccepted(url.CreateAsa(testModel.BaseUrl), doneTransaction)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
				internalTesting.MockGetOk(url.ReadSpecificDevice(testModel.BaseUrl, readOutput.Uid), readSpecificOutput)
			},
			assertFunc: func(actualOutput *asa.ReadOutput, actualSpecificDeviceOutput *asa.ReadSpecificOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, actualOutput)
				assert.Nil(t, actualSpecificDeviceOutput)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, "ASDM version mismatch.")
			},
		},
		{
			testName: "fails onboards ASA if transaction fails",
			input:    createInput,

			setupFunc: func(input asa.CreateInput) {
				internalTesting.MockGetOk(url.ReadConnectorByUid(testModel.BaseUrl, cdgReadOutput.Uid), cdgReadOutput)
				internalTesting.MockPostError(url.CreateAsa(testModel.BaseUrl), errorTransaction)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
			},

			assertFunc: func(actualOutput *asa.ReadOutput, actualSpecificDeviceOutput *asa.ReadSpecificOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, actualOutput)
				assert.Nil(t, actualSpecificDeviceOutput)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, errorTransaction.ErrorMessage)
			},
		},
		{
			testName: "fails onboards ASA if trigger transaction fails",
			input:    createInput,

			setupFunc: func(input asa.CreateInput) {
				internalTesting.MockPostError(url.CreateAsa(testModel.BaseUrl), "post error")
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, readOutput.Uid), readOutput)
			},

			assertFunc: func(actualOutput *asa.ReadOutput, actualSpecificDeviceOutput *asa.ReadSpecificOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, actualOutput)
				assert.Nil(t, actualSpecificDeviceOutput)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, "post error")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, specificDeviceOuput, err := asa.Create(
				context.Background(),
				*http.MustNewWithConfig(testModel.BaseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, specificDeviceOuput, err, t)
		})
	}
}
