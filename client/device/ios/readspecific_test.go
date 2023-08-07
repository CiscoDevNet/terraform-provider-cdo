package ios

import (
	"context"
	"reflect"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/device/ios/iosconfig"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestIosReadSpecific(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	deviceUid := "00000000-0000-0000-0000-000000000000"

	specificDevice := ReadSpecificOutput{
		SpecificUid: "11111111-1111-1111-1111-111111111111",
		State:       iosconfig.IosConfigStateDone,
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
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatalf("output is nil!")
				}

				if !reflect.DeepEqual(specificDevice, *output) {
					t.Errorf("expected: %+v, got: %+v", specificDevice, *output)
				}
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
				if err == nil {
					t.Error("error is nil!")
				}

				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", *output)
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := ReadSpecific(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
