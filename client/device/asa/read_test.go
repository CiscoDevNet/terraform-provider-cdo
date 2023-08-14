package asa_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestAsaRead(t *testing.T) {
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
		input      asa.ReadInput
		setupFunc  func()
		assertFunc func(output *asa.ReadOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully reads ASA",
			input: asa.ReadInput{
				Uid: asaDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceReadToRespondSuccessfully(asaDevice)
			},

			assertFunc: func(output *asa.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedReadOutput := asa.ReadOutput{
					Uid:             asaDevice.Uid,
					Name:            asaDevice.Name,
					CreatedDate:     asaDevice.CreatedDate,
					LastUpdatedDate: asaDevice.LastUpdatedDate,
					DeviceType:      asaDevice.DeviceType,
					LarUid:          asaDevice.LarUid,
					LarType:         asaDevice.LarType,
					Ipv4:            asaDevice.Ipv4,
					Host:            asaDevice.Host,
					Port:            asaDevice.Port,
				}
				assert.Equal(t, expectedReadOutput, *output)
			},
		},

		{
			testName: "returns error when the remote service reading the ASA encounters an issue",
			input: asa.ReadInput{
				Uid: asaDevice.Uid,
			},

			setupFunc: func() {
				configureDeviceReadToRespondWithError(asaDevice.Uid)
			},

			assertFunc: func(output *asa.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := asa.Read(context.Background(), *http.MustNewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
