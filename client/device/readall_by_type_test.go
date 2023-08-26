package device_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestDeviceReadAllByType(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validDevice1 := device.
		NewReadOutputBuilder().
		AsCdfmc().
		WithUid(deviceUid1).
		WithName(deviceName1).
		Build()

	validDevice2 := device.
		NewReadOutputBuilder().
		AsCdfmc().
		WithUid(deviceUid2).
		WithName(deviceName2).
		Build()

	validReadAllOutput := device.ReadAllByTypeOutput{
		validDevice1,
		validDevice2,
	}

	testCases := []struct {
		testName   string
		targetType devicetype.Type
		setupFunc  func()
		assertFunc func(output *device.ReadAllByTypeOutput, err error, t *testing.T)
	}{
		{
			testName:   "successfully read devices by type",
			targetType: devicetype.Cdfmc,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAllDevicesByType(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validReadAllOutput),
				)
			},
			assertFunc: func(output *device.ReadAllByTypeOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validReadAllOutput, *output)
			},
		},
		{
			testName:   "return error when read devices by type error",
			targetType: devicetype.Cdfmc,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.CreateDevice(baseUrl),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *device.ReadAllByTypeOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := device.ReadAllByType(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				device.NewReadAllByTypeInput(testCase.targetType),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
