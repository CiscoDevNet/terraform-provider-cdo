package ios

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestIosDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	iosDevice := device.NewReadOutputBuilder().
		AsIos().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-ios").
		OnboardedUsingOnPremConnector("88888888-8888-8888-8888-888888888888").
		WithLocation("10.10.0.1", 443).
		Build()

	testCases := []struct {
		testName   string
		input      DeleteInput
		setupFunc  func()
		assertFunc func(output *DeleteOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully deletes iOS device",
			input: DeleteInput{
				Uid: iosDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceDeleteToRespondSuccessfully(iosDevice.Uid)
			},

			assertFunc: func(output *DeleteOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedDeleteOutput := DeleteOutput{}
				assert.Equal(t, expectedDeleteOutput, *output)
			},
		},

		{
			testName: "returns error when the remote service deleting the iOS device encounters an issue",
			input: DeleteInput{
				Uid: iosDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceDeleteToRespondWithError(iosDevice.Uid)
			},

			assertFunc: func(output *DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := Delete(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
