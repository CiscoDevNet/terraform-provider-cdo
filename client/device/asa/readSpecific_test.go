package asa_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/device/asa"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestAsaReadSpecific(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	deviceUid := "00000000-0000-0000-0000-000000000000"

	specificDevice := asa.NewReadSpecificOutputBuilder().
		WithSpecificUid("11111111-1111-1111-1111-111111111111").
		InDoneState().
		Build()

	testCases := []struct {
		testName   string
		input      asa.ReadSpecificInput
		setupFunc  func()
		assertFunc func(output *asa.ReadSpecificOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully reads ASA specific device",
			input: asa.ReadSpecificInput{
				Uid: deviceUid,
			},

			setupFunc: func() {
				configureDeviceReadSpecificToRespondSuccessfully(deviceUid, device.ReadSpecificOutput(specificDevice))
			},

			assertFunc: func(output *asa.ReadSpecificOutput, err error, t *testing.T) {
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
			testName: "returns error when the remote service reading the ASA specific device encounters an issue",
			input: asa.ReadSpecificInput{
				Uid: deviceUid,
			},

			setupFunc: func() {
				configureDeviceReadSpecificToRespondWithError(deviceUid)
			},

			assertFunc: func(output *asa.ReadSpecificOutput, err error, t *testing.T) {
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

			output, err := asa.ReadSpecific(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
