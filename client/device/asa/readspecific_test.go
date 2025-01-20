package asa_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
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
				configureDeviceReadSpecificToRespondSuccessfully(deviceUid, asa.ReadSpecificOutput(specificDevice))
			},

			assertFunc: func(output *asa.ReadSpecificOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, specificDevice, *output)
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
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := asa.ReadSpecific(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
