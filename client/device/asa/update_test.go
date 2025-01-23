package asa_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestAsaUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	onPremConnector := connector.NewConnectorOutputBuilder().
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
			testName: "fail to upgrade an ASA device if version not compatible",
			input: asa.UpdateInput{
				Uid:             asaDevice.Uid,
				SoftwareVersion: "9.8(4)",
				AsdmVersion:     "7.16(1)",
			},
			setupFunc: func(input asa.UpdateInput) {
				configureCompatibleVersionsToRespondSuccessfully(input.Uid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: "9.16(4)", AsdmVersion: "7.12(1)"},
						{SoftwareVersion: "9.16(6)100", AsdmVersion: "7.12(2)"},
					},
					Count: 2,
				})
			},
			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "fail to upgrade an ASA device if fetching version compatibility matrix fails",
			input: asa.UpdateInput{
				Uid:             asaDevice.Uid,
				SoftwareVersion: "9.8(4)",
				AsdmVersion:     "7.16(1)",
			},
			setupFunc: func(input asa.UpdateInput) {
				configureCompatibleVersionsToFailToRespond(input.Uid)
			},
			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "fail to upgrade an ASA device if triggering upgrade fails",
			input: asa.UpdateInput{
				Uid:             asaDevice.Uid,
				SoftwareVersion: "9.8(4)",
				AsdmVersion:     "7.16(1)",
			},
			setupFunc: func(input asa.UpdateInput) {
				configureCompatibleVersionsToRespondSuccessfully(input.Uid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: "9.16(4)", AsdmVersion: "7.12(1)"},
						{SoftwareVersion: "9.16(6)100", AsdmVersion: "7.12(2)"},
					},
					Count: 2,
				})
				configureUpgradeAsaToFailToRespond(input.Uid)
			},
			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "successfully upgrade an ASA device",
			input: asa.UpdateInput{
				Uid:             asaDevice.Uid,
				SoftwareVersion: "9.8(4)",
				AsdmVersion:     "7.16(1)",
			},
			setupFunc: func(input asa.UpdateInput) {
				transactionUid := uuid.New().String()
				inProgressTransaction := transaction.Type{
					TransactionUid:  uuid.New().String(),
					TenantUid:       uuid.New().String(),
					EntityUid:       uuid.New().String(),
					EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/inventory/devices/" + input.Uid,
					PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
					SubmissionTime:  "2024-09-10T20:10:00Z",
					LastUpdatedTime: "2024-10-10T20:10:00Z",
					Type:            transactiontype.UPGRADE_ASA,
					Status:          transactionstatus.IN_PROGRESS,
				}
				doneTransaction := transaction.Type{
					TransactionUid:  inProgressTransaction.TransactionUid,
					TenantUid:       inProgressTransaction.TenantUid,
					EntityUid:       inProgressTransaction.EntityUid,
					EntityUrl:       inProgressTransaction.EntityUrl,
					PollingUrl:      inProgressTransaction.PollingUrl,
					SubmissionTime:  inProgressTransaction.SubmissionTime,
					LastUpdatedTime: "2024-10-10T20:11:00Z",
					Type:            inProgressTransaction.Type,
					Status:          transactionstatus.DONE,
				}
				configureCompatibleVersionsToRespondSuccessfully(input.Uid, model.CdoListResponse[asa.CompatibleVersion]{
					Items: []asa.CompatibleVersion{
						{SoftwareVersion: input.SoftwareVersion, AsdmVersion: input.AsdmVersion},
						{SoftwareVersion: "9.16(6)100", AsdmVersion: "7.12(2)"},
					},
					Count: 2,
				})
				configureUpgradeAsaToRespondSuccessfully(input.Uid, inProgressTransaction)
				configureDeviceUpdateToRespondSuccessfully(input.Uid, asaDevice)
				internalTesting.MockGetOk(fmt.Sprintf("%s/api/rest/v1/transactions/%s", baseUrl, transactionUid), doneTransaction)
				configureDeviceUpdateToRespondSuccessfully(input.Uid, asaDevice)
				configureDeviceReadToRespondSuccessfully(device.ReadOutput{Uid: input.Uid, State: "DONE", Status: "IDLE", ConnectivityState: 1})
			},
			assertFunc: func(input asa.UpdateInput, output *asa.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
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
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, asaConfig)
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
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, asaConfig)
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureConnectorReadToRespondSuccessfully(onPremConnector)

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

				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, asaConfig)
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

				configureConnectorReadToRespondSuccessfully(onPremConnector)

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
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, asaConfig)
				configureDeviceReadToRespondWithError(asaDeviceOnboardedByOnPremConnector.Uid)

				configureConnectorReadToRespondSuccessfully(onPremConnector)

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
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, asaConfig)
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureConnectorReadToRespondWithError(onPremConnector.Uid)

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
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, asaConfig)
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureConnectorReadToRespondSuccessfully(onPremConnector)

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
				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, asaConfig)
				configureDeviceReadToRespondSuccessfully(asaDeviceOnboardedByOnPremConnector)

				configureConnectorReadToRespondSuccessfully(onPremConnector)

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

				configureDeviceReadSpecificToRespondSuccessfully(input.Uid, asaConfig)
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
