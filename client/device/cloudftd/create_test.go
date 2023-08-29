package cloudftd_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCreateCloudFtd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName   string
		input      cloudftd.CreateInput
		setupFunc  func()
		assertFunc func(output *cloudftd.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully create Cloud FTD",
			setupFunc: func() {
				readFmcIsSuccessful(true)
				readFmcDomainInfoIsSuccessful(true)
				readFmcAccessPoliciesIsSuccessful(true)
				createFtdIsSuccessful(true)
				readFtdSpecificDeviceIsSuccessful(true)
				triggerFtdOnboardingIsSuccessful(true)
				generateFtdConfigureManagerCommandIsSuccessful(true)
			},
			assertFunc: func(output *cloudftd.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validFtdCreateOutput, *output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := cloudftd.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
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

func readFmcAccessPoliciesIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadAccessPolicies(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadAccessPoliciesOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadAccessPolicies(baseUrl, fmcDomainUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func createFtdIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodPost,
			url.CreateDevice(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validCreateFtdOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodPost,
			url.CreateDevice(baseUrl),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func readFtdSpecificDeviceIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadSpecificDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadFtdSpecificDeviceOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadSpecificDevice(baseUrl, ftdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func triggerFtdOnboardingIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodPut,
			url.UpdateSpecificCloudFtd(baseUrl, createSpecificFtdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validUpdateSpecificFtdOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodPut,
			url.UpdateSpecificCloudFtd(baseUrl, createSpecificFtdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}

func generateFtdConfigureManagerCommandIsSuccessful(success bool) {
	if success {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadDevice(baseUrl, createdFtdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadFtdGeneratedCommandOutput),
		)
	} else {
		httpmock.RegisterResponder(
			http.MethodGet,
			url.ReadDevice(baseUrl, createdFtdUid),
			httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
		)
	}
}
