package asa_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestAsaUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	onPremConnector := sdc.NewSdcOutputBuilder().
		AsOnPremConnector().
		WithUid("00000000-0000-0000-0000-000000000000").
		WithName("MyCloudConnector").
		WithTenantUid("66666666-6666-6666-6666-6666666666666").
		Build()

	asaDevice := device.NewReadOutputBuilder().
		AsAsa().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-asa").
		OnboardedUsingCloudConnector("88888888-8888-8888-8888-888888888888").
		WithLocation("10.10.0.1", 443).
		Build()

	asaDeviceOnboardedByOnPremConnector := device.NewReadOutputBuilder().
		AsAsa().
		WithUid("33333333-3333-3333-3333-333333333333").
		WithName("my-asa").
		OnboardedUsingOnPremConnector(onPremConnector.Uid).
		WithLocation("10.10.0.1", 443).
		Build()

	asaConfig := asa.NewReadSpecificOutputBuilder().
		WithSpecificUid("22222222-2222-2222-2222-222222222222").
		InDoneState().
		Build()

	testCases := []struct {
		testName   string
		input      asa.UpdateInput
		setupFunc  func(input asa.UpdateInput)
		assertFunc func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully updates ASA name",
			input: asa.UpdateInput{
				Uid:  asaDevice.Uid,
				Name: "new-name",
			},

			setupFunc: func(input asa.UpdateInput) {
				updatedDevice := asaDevice
				updatedDevice.Name = input.Name
				configureDeviceUpdateToRespondSuccessfully(input.Uid, updatedDevice)
				configureAsaConfigReadToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.ReadOutput{Uid: asaConfig.SpecificUid, State: state.DONE})
				configureDeviceReadToRespondSuccessfully(device.ReadOutput{Uid: input.Uid, State: "DONE", Status: "IDLE", ConnectivityState: 1})
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedUpdateOutput := asaDevice
				expectedUpdateOutput.Name = input.Name
				assert.Equal(t, expectedUpdateOutput, *output)
			},
		},

		{
			testName: "successfully updates ASA credentials",
			input: asa.UpdateInput{
				Uid:      asaDevice.Uid,
				Username: "lockhart",
				Password: "not a valid password",
			},

			setupFunc: func(input asa.UpdateInput) {
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, device.ReadSpecificOutput(asaConfig))
				configureDeviceReadToRespondSuccessfully(asaDevice)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.UpdateOutput{Uid: asaConfig.SpecificUid})
				configureAsaConfigReadToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.ReadOutput{Uid: asaConfig.SpecificUid, State: state.DONE})
				configureDeviceUpdateToRespondSuccessfully(input.Uid, asaDevice)
				configureDeviceReadToRespondSuccessfully(device.ReadOutput{Uid: input.Uid, State: "DONE", Status: "IDLE", ConnectivityState: 1})
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, asaDevice, *output)
			},
		},

		{
			testName: "successfully updates ASA credentials via an OnPrem Connector",
			input: asa.UpdateInput{
				Uid:      asaDeviceOnboardedByOnPremConnector.Uid,
				Username: "lockhart",
				Password: "not a valid password",
			},

			setupFunc: func(input asa.UpdateInput) {
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, device.ReadSpecificOutput(asaConfig))
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureSdcReadToRespondSuccessfully(onPremConnector)

				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.UpdateOutput{Uid: asaConfig.SpecificUid})
				configureAsaConfigReadToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.ReadOutput{Uid: asaConfig.SpecificUid, State: state.DONE})
				configureDeviceUpdateToRespondSuccessfully(input.Uid, asaDeviceOnboardedByOnPremConnector)
				configureDeviceReadToRespondSuccessfully(device.ReadOutput{Uid: input.Uid, State: "DONE", Status: "IDLE", ConnectivityState: 1})
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, asaDeviceOnboardedByOnPremConnector, *output)
			},
		},

		{
			testName: "successfully updates ASA location",
			input: asa.UpdateInput{
				Uid:      asaDevice.Uid,
				Location: "10.10.5.4:443",
			},

			setupFunc: func(input asa.UpdateInput) {
				updatedDevice := asaDevice
				updatedDevice.Host = "10.10.5.4"
				updatedDevice.Port = "443"

				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, device.ReadSpecificOutput(asaConfig))
				configureDeviceReadToRespondSuccessfully(asaDevice)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.UpdateOutput{Uid: asaConfig.SpecificUid})
				configureAsaConfigReadToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.ReadOutput{Uid: asaConfig.SpecificUid, State: state.DONE})
				configureDeviceUpdateToRespondSuccessfully(input.Uid, updatedDevice)
				configureDeviceReadToRespondSuccessfully(device.ReadOutput{Uid: input.Uid, State: "DONE", Status: "IDLE", ConnectivityState: 1})
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				updatedDevice := asaDevice
				updatedDevice.Host = "10.10.5.4"
				updatedDevice.Port = "443"
				assert.Equal(t, updatedDevice, *output)
			},
		},

		{
			testName: "returns error when device read specific call encounters an issue",
			input: asa.UpdateInput{
				Uid:      asaDeviceOnboardedByOnPremConnector.Uid,
				Username: "lockhart",
				Password: "not a valid password",
			},

			setupFunc: func(input asa.UpdateInput) {
				configureDeviceReadSpecificToRespondWithError(input.Uid)
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureSdcReadToRespondSuccessfully(onPremConnector)

				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.UpdateOutput{Uid: asaConfig.SpecificUid})
				configureDeviceUpdateToRespondSuccessfully(input.Uid, asaDeviceOnboardedByOnPremConnector)
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},

		{
			testName: "returns error when device read call encounters an issue",
			input: asa.UpdateInput{
				Uid:      asaDeviceOnboardedByOnPremConnector.Uid,
				Username: "lockhart",
				Password: "not a valid password",
			},

			setupFunc: func(input asa.UpdateInput) {
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, device.ReadSpecificOutput(asaConfig))
				configureDeviceReadToRespondWithError(asaDeviceOnboardedByOnPremConnector.Uid)

				configureSdcReadToRespondSuccessfully(onPremConnector)

				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.UpdateOutput{Uid: asaConfig.SpecificUid})
				configureDeviceUpdateToRespondSuccessfully(input.Uid, asaDeviceOnboardedByOnPremConnector)
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},

		{
			testName: "returns error when sdc read call encounters an issue",
			input: asa.UpdateInput{
				Uid:      asaDeviceOnboardedByOnPremConnector.Uid,
				Username: "lockhart",
				Password: "not a valid password",
			},

			setupFunc: func(input asa.UpdateInput) {
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, device.ReadSpecificOutput(asaConfig))
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureSdcReadToRespondWithError(onPremConnector.Uid)

				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.UpdateOutput{Uid: asaConfig.SpecificUid})
				configureDeviceUpdateToRespondSuccessfully(input.Uid, asaDeviceOnboardedByOnPremConnector)
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},

		{
			testName: "returns error when asa config update call encounters an issue",
			input: asa.UpdateInput{
				Uid:      asaDeviceOnboardedByOnPremConnector.Uid,
				Username: "lockhart",
				Password: "not a valid password",
			},

			setupFunc: func(input asa.UpdateInput) {
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, device.ReadSpecificOutput(asaConfig))
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureSdcReadToRespondSuccessfully(onPremConnector)

				configureAsaConfigUpdateToRespondWithError(asaConfig.SpecificUid)
				configureDeviceUpdateToRespondSuccessfully(input.Uid, asaDeviceOnboardedByOnPremConnector)
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},

		{
			testName: "returns error when device update call encounters an issue",
			input: asa.UpdateInput{
				Uid:      asaDeviceOnboardedByOnPremConnector.Uid,
				Username: "lockhart",
				Password: "not a valid password",
			},

			setupFunc: func(input asa.UpdateInput) {
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, device.ReadSpecificOutput(asaConfig))
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureSdcReadToRespondSuccessfully(onPremConnector)

				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.UpdateOutput{Uid: asaConfig.SpecificUid})
				configureDeviceUpdateToRespondWithError(asaDeviceOnboardedByOnPremConnector.Uid)
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)

				assert.NotNil(t, err)
			},
		},

		{
			testName: "returns error when asa config read call encounters an issue",
			input: asa.UpdateInput{
				Uid:      asaDevice.Uid,
				Location: "10.10.5.4:443",
			},

			setupFunc: func(input asa.UpdateInput) {
				updatedDevice := asaDevice
				updatedDevice.Host = "10.10.5.4"
				updatedDevice.Port = "443"

				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, device.ReadSpecificOutput(asaConfig))
				configureDeviceReadToRespondSuccessfully(asaDevice)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.SpecificUid, asaconfig.UpdateOutput{Uid: asaConfig.SpecificUid})
				configureAsaConfigReadToRespondWithError(asaConfig.SpecificUid)
				configureDeviceUpdateToRespondSuccessfully(input.Uid, updatedDevice)
			},

			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := asa.Update(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(testCase.input, output, err, t)
		})
	}
}
