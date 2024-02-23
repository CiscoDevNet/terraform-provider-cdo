package ios_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
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
		WithTags(tags.New("tags1", "tags2", "tags3")).
		Build()

	testCases := []struct {
		testName   string
		input      ios.UpdateInput
		setupFunc  func(input ios.UpdateInput)
		assertFunc func(input ios.UpdateInput, output *ios.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully updates iOS name",
			input: ios.UpdateInput{
				Uid:  iosDevice.Uid,
				Name: "new-name",
			},

			setupFunc: func(input ios.UpdateInput) {
				updatedDevice := iosDevice
				updatedDevice.Name = input.Name
				configureDeviceUpdateToRespondSuccessfully(updatedDevice)
			},

			assertFunc: func(input ios.UpdateInput, output *ios.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedUpdateOutput := iosDevice
				expectedUpdateOutput.Name = input.Name
				assert.Equal(t, expectedUpdateOutput, *output)
			},
		},

		{
			testName: "returns error when device update call encounters an issue",
			input: ios.UpdateInput{
				Uid:  iosDevice.Uid,
				Name: "new-name",
			},

			setupFunc: func(input ios.UpdateInput) {
				configureDeviceUpdateToRespondWithError(iosDevice.Uid)
			},

			assertFunc: func(input ios.UpdateInput, output *ios.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := ios.Update(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(testCase.input, output, err, t)
		})
	}
}
