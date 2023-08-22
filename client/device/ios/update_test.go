package ios

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestIosUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	onPremConnector := connector.NewConnectorOutputBuilder().
		AsOnPremConnector().
		WithUid("00000000-0000-0000-0000-000000000000").
		WithName("MyOnPremConnector").
		WithTenantUid("66666666-6666-6666-6666-6666666666666").
		Build()

	iosDevice := device.NewReadOutputBuilder().
		AsIos().
		WithUid("33333333-3333-3333-3333-333333333333").
		WithName("my-ios").
		OnboardedUsingOnPremConnector(onPremConnector.Uid).
		WithLocation("10.10.0.1", 443).
		Build()

	testCases := []struct {
		testName   string
		input      UpdateInput
		setupFunc  func(input UpdateInput)
		assertFunc func(input UpdateInput, output *UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully updates iOS name",
			input: UpdateInput{
				Uid:  iosDevice.Uid,
				Name: "new-name",
			},

			setupFunc: func(input UpdateInput) {
				updatedDevice := iosDevice
				updatedDevice.Name = input.Name
				configureDeviceUpdateToRespondSuccessfully(updatedDevice)
			},

			assertFunc: func(input UpdateInput, output *UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedUpdateOutput := iosDevice
				expectedUpdateOutput.Name = input.Name
				assert.Equal(t, expectedUpdateOutput, *output)
			},
		},

		{
			testName: "returns error when device update call encounters an issue",
			input: UpdateInput{
				Uid:  iosDevice.Uid,
				Name: "new-name",
			},

			setupFunc: func(input UpdateInput) {
				configureDeviceUpdateToRespondWithError(iosDevice.Uid)
			},

			assertFunc: func(input UpdateInput, output *UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := Update(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(testCase.input, output, err, t)
		})
	}
}
