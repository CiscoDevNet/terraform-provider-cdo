package ios

import (
	"context"
	"reflect"
	"testing"

	"github.com/CiscoDevnet/go-client/device"
	"github.com/CiscoDevnet/go-client/internal/http"
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
		Build()

	testCases := []struct {
		testName   string
		input      ReadInput
		setupFunc  func()
		assertFunc func(output *ReadOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully reads iOS",
			input: ReadInput{
				Uid: iosDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceReadToRespondSuccessfully(iosDevice)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatalf("output is nil!")
				}

				expectedReadOutput := ReadOutput{
					Uid:             iosDevice.Uid,
					Name:            iosDevice.Name,
					CreatedDate:     iosDevice.CreatedDate,
					LastUpdatedDate: iosDevice.LastUpdatedDate,
					DeviceType:      iosDevice.DeviceType,
					LarUid:          iosDevice.LarUid,
					LarType:         iosDevice.LarType,
					Ipv4:            iosDevice.Ipv4,
					Host:            iosDevice.Host,
					Port:            iosDevice.Port,
				}
				if !reflect.DeepEqual(expectedReadOutput, *output) {
					t.Errorf("expected: %+v, got: %+v", iosDevice, *output)
				}
			},
		},

		{
			testName: "returns error when the remote service reading the iOS encounters an issue",
			input: ReadInput{
				Uid: iosDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceReadToRespondWithError(iosDevice.Uid)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
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

			output, err := Read(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
