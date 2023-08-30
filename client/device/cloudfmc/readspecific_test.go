package cloudfmc_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestReadSpecificCloudFmc(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validReadSpecificOutput := cloudfmc.NewReadSpecificOutputBuilder().
		SpecificUid(specificDeviceUid).
		DomainUid(domainUid).
		Status(status).
		State(deviceState).
		Build()

	testCases := []struct {
		testName   string
		setupFunc  func()
		assertFunc func(output *cloudfmc.ReadSpecificOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully read Cloud FMC specific device",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadSpecificDevice(baseUrl, fmcUid),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadSpecificOutput),
				)
			},
			assertFunc: func(output *cloudfmc.ReadSpecificOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validReadSpecificOutput, *output)
			},
		},
		{
			testName: "error when read Cloud FMC specific device error",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadSpecificDevice(baseUrl, fmcUid),
					httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *cloudfmc.ReadSpecificOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := cloudfmc.ReadSpecific(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				cloudfmc.NewReadSpecificInput(fmcUid),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
