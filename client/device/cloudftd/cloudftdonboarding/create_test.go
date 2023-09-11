package cloudftdonboarding_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd/cloudftdonboarding"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
	"time"
)

func TestCloudFtdOnboardingCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName   string
		input      cloudftdonboarding.CreateInput
		setupFunc  func()
		assertFunc func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successful ftd onboarding",
			input:    cloudftdonboarding.NewCreateInput("unittest-device-uid"),
			setupFunc: func() {
				ReadFmcIsSuccessful(true)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := cloudftdonboarding.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig("https://unittest.cdo.cisco.com", "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func ReadFmcIsSuccessful(success bool) {
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
