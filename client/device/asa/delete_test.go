package asa_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsaDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	asaDevice := device.NewReadOutputBuilder().
		AsAsa().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-asa").
		OnboardedUsingCloudConnector("88888888-8888-8888-8888-888888888888").
		WithLocation("10.10.0.1", 443).
		Build()

	testCases := []struct {
		testName   string
		input      asa.DeleteInput
		setupFunc  func()
		assertFunc func(output *asa.DeleteOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully deletes ASA",
			input: asa.DeleteInput{
				Uid: asaDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceDeleteToRespondSuccessfully(asaDevice.Uid)
			},

			assertFunc: func(output *asa.DeleteOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedDeleteOutput := asa.DeleteOutput{}
				assert.Equal(t, expectedDeleteOutput, *output)
			},
		},

		{
			testName: "returns error when the remote service deleting the ASA encounters an issue",
			input: asa.DeleteInput{
				Uid: asaDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceDeleteToRespondWithError(asaDevice.Uid)
			},

			assertFunc: func(output *asa.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := asa.Delete(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
