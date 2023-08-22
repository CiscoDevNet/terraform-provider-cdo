package ios

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestIosReadSpecific(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	deviceUid := "00000000-0000-0000-0000-000000000000"

	specificDevice := ReadSpecificOutput{
		SpecificUid: "11111111-1111-1111-1111-111111111111",
		State:       state.DONE,
		Namespace:   "targets",
		Type:        "device",
	}

	testCases := []struct {
		testName   string
		input      ReadSpecificInput
		setupFunc  func()
		assertFunc func(output *ReadSpecificOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully reads iOS specific device",
			input: ReadSpecificInput{
				Uid: deviceUid,
			},

			setupFunc: func() {
				configureDeviceReadSpecificToRespondSuccessfully(deviceUid, device.ReadSpecificOutput(specificDevice))
			},

			assertFunc: func(output *ReadSpecificOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, specificDevice, *output)
			},
		},

		{
			testName: "returns error when the remote service reading the iOS specific device encounters an issue",
			input: ReadSpecificInput{
				Uid: deviceUid,
			},

			setupFunc: func() {
				configureDeviceReadSpecificToRespondWithError(deviceUid)
			},

			assertFunc: func(output *ReadSpecificOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := ReadSpecific(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
