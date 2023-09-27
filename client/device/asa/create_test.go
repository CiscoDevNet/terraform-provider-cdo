package asa_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/featureflag"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth/role"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/jarcoal/httpmock"
)

func TestAsaCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	asaDevice := device.NewReadOutputBuilder().
		AsAsa().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-asa").
		OnboardedUsingCloudConnector("88888888-8888-8888-8888-888888888888").
		WithLocation("10.10.0.1", 443).
		WithTags(tags.New("tag1", "tag2")).
		Build()

	asaDeviceUsingSdc := device.NewReadOutputBuilder().
		AsAsa().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-asa").
		OnboardedUsingOnPremConnector("99999999-9999-9999-9999-999999999999").
		WithLocation("10.10.0.1", 443).
		Build()

	asaSpecificDevice := device.ReadSpecificOutput{
		SpecificUid: "22222222-2222-2222-2222-222222222222",
		State:       "DONE",
		Namespace:   "devices",
		Type:        "asa",
	}
	asaConfig := asaconfig.ReadOutput{
		Uid:   asaSpecificDevice.SpecificUid,
		State: state.DONE,
	}
	readApiTokenInfo_NoNewModelFeatureFlag := auth.Info{UserAuthentication: auth.Authentication{
		Authorities: []auth.Authority{
			{Authority: role.Admin},
		},
		Details: auth.Details{
			TenantUid:              "11111111-1111-1111-1111-111111111111",
			TenantName:             "",
			SseTenantUid:           "",
			TenantOrganizationName: "",
			TenantDbFeatures:       "{}",
			TenantUserRoles:        "",
			TenantDatabaseName:     "",
			TenantPayType:          "",
		},
		Authenticated: false,
		Principle:     "",
		Name:          "",
	}}

	readApiTokenInfo_NewModelFeatureFlag := auth.Info{UserAuthentication: auth.Authentication{
		Authorities: []auth.Authority{
			{Authority: role.Admin},
		},
		Details: auth.Details{
			TenantUid:              "11111111-1111-1111-1111-111111111111",
			TenantName:             "",
			SseTenantUid:           "",
			TenantOrganizationName: "",
			TenantDbFeatures:       fmt.Sprintf("{\"%s\":true}", featureflag.AsaConfigurationObjectMigration),
			TenantUserRoles:        "",
			TenantDatabaseName:     "",
			TenantPayType:          "",
		},
		Authenticated: false,
		Principle:     "",
		Name:          "",
	}}

	validConnector := connector.NewConnectorOutputBuilder().
		WithName("CloudDeviceGateway").
		WithUid(asaDeviceUsingSdc.ConnectorUid).
		WithTenantUid("44444444-4444-4444-4444-444444444444").
		AsOnPremConnector().
		Build()

	testCases := []struct {
		testName   string
		input      asa.CreateInput
		setupFunc  func(input asa.CreateInput)
		assertFunc func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T)
	}{
		{
			testName: "successfully onboards ASA when using CDG",
			input: asa.CreateInput{
				Name:              asaDevice.Name,
				ConnectorType:     asaDevice.ConnectorType,
				SocketAddress:     asaDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoSuccessfully(readApiTokenInfo_NoNewModelFeatureFlag)
				configureDeviceCreateToRespondSuccessfully(asaDevice)
				configureDeviceReadSpecificToRespondSuccessfully(asaDevice.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedCreatedOutput := asa.CreateOutput{
					Uid:           asaDevice.Uid,
					Name:          asaDevice.Name,
					DeviceType:    asaDevice.DeviceType,
					Host:          asaDevice.Host,
					Port:          asaDevice.Port,
					SocketAddress: asaDevice.SocketAddress,
					ConnectorType: asaDevice.ConnectorType,
					ConnectorUid:  asaDevice.ConnectorUid,
					Tags:          asaDevice.Tags,
				}
				assert.Equal(t, expectedCreatedOutput, *output)

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDevice.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "returns error when invalid api token is given",
			input: asa.CreateInput{
				Name:              asaDevice.Name,
				ConnectorType:     asaDevice.ConnectorType,
				SocketAddress:     asaDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoFailed()
				configureDeviceCreateToRespondSuccessfully(asaDevice)
				configureDeviceReadSpecificToRespondSuccessfully(asaDevice.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "should call create device with new policy model metadata if feature flag for new model is enable",
			input: asa.CreateInput{
				Name:              asaDevice.Name,
				ConnectorType:     asaDevice.ConnectorType,
				SocketAddress:     asaDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoSuccessfully(readApiTokenInfo_NewModelFeatureFlag)
				configureDeviceCreateToRespondSuccessfullyWithNewModel(t, asaDevice)
				configureDeviceReadSpecificToRespondSuccessfully(asaDevice.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
			},
		},

		{
			testName: "should not call create device with new policy model metadata if feature flag for new model is not enabled",
			input: asa.CreateInput{
				Name:              asaDevice.Name,
				ConnectorType:     asaDevice.ConnectorType,
				SocketAddress:     asaDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoSuccessfully(readApiTokenInfo_NoNewModelFeatureFlag)
				configureDeviceCreateToRespondSuccessfullyWithoutNewModel(t, asaDevice)
				configureDeviceReadSpecificToRespondSuccessfully(asaDevice.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
			},
		},

		{
			testName: "successfully onboards ASA when using CDG after recovering from certificate error",

			input: asa.CreateInput{
				Name:              asaDevice.Name,
				ConnectorType:     asaDevice.ConnectorType,
				SocketAddress:     asaDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoSuccessfully(readApiTokenInfo_NoNewModelFeatureFlag)
				configureDeviceCreateToRespondSuccessfully(asaDevice)
				configureDeviceReadSpecificToRespondSuccessfully(asaDevice.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: state.ERROR,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureDeviceUpdateToRespondSuccessfully(asaDevice.Uid, asaDevice)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedCreatedOutput := asa.CreateOutput{
					Uid:           asaDevice.Uid,
					Name:          asaDevice.Name,
					DeviceType:    asaDevice.DeviceType,
					Host:          asaDevice.Host,
					Port:          asaDevice.Port,
					SocketAddress: asaDevice.SocketAddress,
					ConnectorType: asaDevice.ConnectorType,
					ConnectorUid:  asaDevice.ConnectorUid,
					Tags:          asaDevice.Tags,
				}
				assert.Equal(t, expectedCreatedOutput, *output)

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDevice.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertDeviceUpdateWasCalledOnce(asaDevice.Uid, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "successfully onboards ASA when using connector",
			input: asa.CreateInput{
				Name:              asaDeviceUsingSdc.Name,
				ConnectorType:     asaDeviceUsingSdc.ConnectorType,
				ConnectorUid:      asaDeviceUsingSdc.ConnectorUid,
				SocketAddress:     asaDeviceUsingSdc.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoSuccessfully(readApiTokenInfo_NoNewModelFeatureFlag)
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureConnectorReadToRespondSuccessfully(validConnector)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedCreatedOutput := asa.CreateOutput{
					Uid:           asaDeviceUsingSdc.Uid,
					Name:          asaDeviceUsingSdc.Name,
					DeviceType:    asaDeviceUsingSdc.DeviceType,
					Host:          asaDeviceUsingSdc.Host,
					Port:          asaDeviceUsingSdc.Port,
					SocketAddress: asaDeviceUsingSdc.SocketAddress,
					ConnectorType: asaDeviceUsingSdc.ConnectorType,
					ConnectorUid:  asaDeviceUsingSdc.ConnectorUid,
				}
				assert.Equal(t, expectedCreatedOutput, *output)

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDeviceUsingSdc.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertConnectorReadByUidWasCalledOnce(validConnector.Uid, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "successfully onboards ASA when using connector after recovering from certificate error",

			input: asa.CreateInput{
				Name:              asaDeviceUsingSdc.Name,
				ConnectorType:     asaDeviceUsingSdc.ConnectorType,
				ConnectorUid:      asaDeviceUsingSdc.ConnectorUid,
				SocketAddress:     asaDeviceUsingSdc.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoSuccessfully(readApiTokenInfo_NoNewModelFeatureFlag)
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: state.ERROR,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureConnectorReadToRespondSuccessfully(validConnector)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedCreatedOutput := asa.CreateOutput{
					Uid:           asaDeviceUsingSdc.Uid,
					Name:          asaDeviceUsingSdc.Name,
					DeviceType:    asaDeviceUsingSdc.DeviceType,
					Host:          asaDeviceUsingSdc.Host,
					Port:          asaDeviceUsingSdc.Port,
					SocketAddress: asaDeviceUsingSdc.SocketAddress,
					ConnectorType: asaDeviceUsingSdc.ConnectorType,
					ConnectorUid:  asaDeviceUsingSdc.ConnectorUid,
				}
				assert.Equal(t, expectedCreatedOutput, *output)

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDeviceUsingSdc.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertDeviceUpdateWasCalledOnce(asaDeviceUsingSdc.Uid, t)
				assertConnectorReadByUidWasCalledOnce(validConnector.Uid, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "returns error when onboarding ASA and create device call experiences issues",

			input: asa.CreateInput{
				Name:              asaDeviceUsingSdc.Name,
				ConnectorType:     asaDeviceUsingSdc.ConnectorType,
				ConnectorUid:      asaDeviceUsingSdc.ConnectorUid,
				SocketAddress:     asaDeviceUsingSdc.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoSuccessfully(readApiTokenInfo_NoNewModelFeatureFlag)
				configureDeviceCreateToRespondWithError()
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: state.ERROR,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureConnectorReadToRespondSuccessfully(validConnector)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when onboarding ASA and read specific device call experiences issues",

			input: asa.CreateInput{
				Name:              asaDeviceUsingSdc.Name,
				ConnectorType:     asaDeviceUsingSdc.ConnectorType,
				ConnectorUid:      asaDeviceUsingSdc.ConnectorUid,
				SocketAddress:     asaDeviceUsingSdc.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureReadApiTokenInfoSuccessfully(readApiTokenInfo_NoNewModelFeatureFlag)
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondWithError(asaDeviceUsingSdc.Uid)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: state.ERROR,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureConnectorReadToRespondSuccessfully(validConnector)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when onboarding ASA and read asa config call experiences issues",

			input: asa.CreateInput{
				Name:              asaDeviceUsingSdc.Name,
				ConnectorType:     asaDeviceUsingSdc.ConnectorType,
				ConnectorUid:      asaDeviceUsingSdc.ConnectorUid,
				SocketAddress:     asaDeviceUsingSdc.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithError(asaConfig.Uid)
				configureConnectorReadToRespondSuccessfully(validConnector)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when onboarding ASA and validConnector read call experiences issues",

			input: asa.CreateInput{
				Name:              asaDeviceUsingSdc.Name,
				ConnectorType:     asaDeviceUsingSdc.ConnectorType,
				ConnectorUid:      asaDeviceUsingSdc.ConnectorUid,
				SocketAddress:     asaDeviceUsingSdc.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: state.ERROR,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureConnectorReadToRespondWithError(validConnector.Uid)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when onboarding ASA and device update call experiences issues",

			input: asa.CreateInput{
				Name:              asaDeviceUsingSdc.Name,
				ConnectorType:     asaDeviceUsingSdc.ConnectorType,
				ConnectorUid:      asaDeviceUsingSdc.ConnectorUid,
				SocketAddress:     asaDeviceUsingSdc.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: state.ERROR,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureConnectorReadToRespondSuccessfully(validConnector)
				configureDeviceUpdateToRespondWithError(asaDeviceUsingSdc.Uid)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},
			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when onboarding ASA and ASA config update call experiences issues",

			input: asa.CreateInput{
				Name:              asaDeviceUsingSdc.Name,
				ConnectorType:     asaDeviceUsingSdc.ConnectorType,
				ConnectorUid:      asaDeviceUsingSdc.ConnectorUid,
				SocketAddress:     asaDeviceUsingSdc.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: state.ERROR,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureConnectorReadToRespondSuccessfully(validConnector)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondWithError(asaConfig.Uid)
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := asa.Create(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
