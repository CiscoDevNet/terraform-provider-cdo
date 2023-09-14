package cloudfmc_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestReadCloudFmc(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validDevice := device.NewReadOutputBuilder().
		AsCloudFmc().
		WithName(deviceName).
		WithUid(deviceUid).
		WithLocation(deviceHost, devicePort).
		WithCreatedDate(deviceCreatedDate).
		WithLastUpdatedDate(deviceLastUpdatedDate).
		OnboardedUsingCloudConnector(deviceCloudConnectorUId).
		Build()

	validReadDeviceOutput := []device.ReadOutput{
		validDevice,
	}
	validReadFmcOutput := validDevice

	testCases := []struct {
		testName   string
		setupFunc  func()
		assertFunc func(output *cloudfmc.ReadOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully read Cloud FMC",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAllDevicesByType(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadDeviceOutput),
				)
			},
			assertFunc: func(output *cloudfmc.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validReadFmcOutput, *output)
			},
		},
		{
			testName: "error when read Cloud FMC error",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAllDevicesByType(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *cloudfmc.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when no Cloud FMC returned",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAllDevicesByType(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, []device.ReadOutput{}),
				)
			},
			assertFunc: func(output *cloudfmc.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "error when multiple Cloud FMCs returned",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAllDevicesByType(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, []device.ReadOutput{validDevice, validDevice}),
				)
			},
			assertFunc: func(output *cloudfmc.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := cloudfmc.Read(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				cloudfmc.NewReadInput(),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
