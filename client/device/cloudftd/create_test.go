package cloudftd_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCloudFtd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testModel := internalTesting.NewRandomModel()

	ftdCreateInput := testModel.FtdCreateInput()
	ftdReadOutput := testModel.FtdReadOutput()
	fmcReadOutput := testModel.FmcReadOutput()
	fmcDomainInfo := testModel.FmcDomainInfo()
	accessPolicy := testModel.ReadAccessPolicies()
	doneTransaction := testModel.CreateDoneTransaction(ftdReadOutput.Uid, transactiontype.CREATE_FTD)
	errorTransaction := testModel.CreateErrorTransaction(ftdReadOutput.Uid, transactiontype.CREATE_FTD)

	testCases := []struct {
		testName   string
		input      cloudftd.CreateInput
		setupFunc  func(createInp cloudftd.CreateInput)
		assertFunc func(output *cloudftd.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successful ftd creation",
			input:    ftdCreateInput,
			setupFunc: func(createInp cloudftd.CreateInput) {
				internalTesting.MockGetOk(url.ReadAllDevicesByType(testModel.BaseUrl), []cloudfmc.ReadOutput{fmcReadOutput})
				internalTesting.MockGetOk(url.ReadFmcDomainInfo(testModel.FmcHost), fmcDomainInfo)
				internalTesting.MockGetOk(url.ReadAccessPolicies(testModel.BaseUrl, testModel.FmcDomainUuid.String()), accessPolicy)
				internalTesting.MockPostAccepted(url.CreateFtd(testModel.BaseUrl), doneTransaction)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, ftdReadOutput.Uid), ftdReadOutput)
			},
			assertFunc: func(output *cloudftd.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, output, cloudftd.FromDeviceReadOutput(&ftdReadOutput))
			},
		},
		{
			testName: "fails ftd creation when transaction fails",
			input:    ftdCreateInput,
			setupFunc: func(createInp cloudftd.CreateInput) {
				internalTesting.MockGetOk(url.ReadAllDevicesByType(testModel.BaseUrl), []cloudfmc.ReadOutput{fmcReadOutput})
				internalTesting.MockGetOk(url.ReadFmcDomainInfo(testModel.FmcHost), fmcDomainInfo)
				internalTesting.MockGetOk(url.ReadAccessPolicies(testModel.BaseUrl, testModel.FmcDomainUuid.String()), accessPolicy)
				internalTesting.MockPostError(url.CreateFtd(testModel.BaseUrl), errorTransaction)
				internalTesting.MockGetOk(url.ReadDevice(testModel.BaseUrl, ftdReadOutput.Uid), ftdReadOutput)
			},
			assertFunc: func(output *cloudftd.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.ErrorContains(t, err, errorTransaction.ErrorMessage)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := cloudftd.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(testModel.BaseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func readFmcIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadAllDevicesByType(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadFmcOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadAllDevicesByType(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func readFmcDomainInfoIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcDomainInfo(fmcHost),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadFmcDomainInfoOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadFmcDomainInfo(fmcHost),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}
