package ios_test

import (
	"context"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios"
	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/jarcoal/httpmock"
)

func TestIosRead(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	iosDevice := device.NewReadOutputBuilder().
		AsIos().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-ios").
		OnboardedUsingOnPremConnector("88888888-8888-8888-8888-888888888888").
		WithLocation("10.10.0.1", 443).
		WithTags(internalTesting.NewTestingTags()).
		Build()

	testCases := []struct {
		testName   string
		input      ios.ReadInput
		setupFunc  func()
		assertFunc func(output *ios.ReadOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully reads iOS",
			input: ios.ReadInput{
				Uid: iosDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceReadToRespondSuccessfully(iosDevice)
			},

			assertFunc: func(output *ios.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedReadOutput := ios.ReadOutput{
					Uid:             iosDevice.Uid,
					Name:            iosDevice.Name,
					CreatedDate:     iosDevice.CreatedDate,
					LastUpdatedDate: iosDevice.LastUpdatedDate,
					DeviceType:      iosDevice.DeviceType,
					ConnectorUid:    iosDevice.ConnectorUid,
					ConnectorType:   iosDevice.ConnectorType,
					SocketAddress:   iosDevice.SocketAddress,
					Host:            iosDevice.Host,
					Port:            iosDevice.Port,
					Tags:            iosDevice.Tags,
				}
				assert.Equal(t, expectedReadOutput, *output)
			},
		},

		{
			testName: "returns error when the remote service reading the iOS encounters an issue",
			input: ios.ReadInput{
				Uid: iosDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceReadToRespondWithError(iosDevice.Uid)
			},

			assertFunc: func(output *ios.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := ios.Read(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
