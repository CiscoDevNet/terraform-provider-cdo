package asa_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/device/asaconfig"
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
		State: asaconfig.AsaConfigStateDone,
	}

	sdc := sdc.NewSdcOutputBuilder().
		WithName("CloudDeviceGateway").
		WithUid(asaDeviceUsingSdc.LarUid).
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
				Name:             asaDevice.Name,
				SdcType:          asaDevice.LarType,
				Ipv4:             asaDevice.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: false,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDevice)
				configureDeviceReadSpecificToRespondSuccessfully(asaDevice.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedCreatedOutput := asa.CreateOutput{
					Uid:        asaDevice.Uid,
					Name:       asaDevice.Name,
					DeviceType: asaDevice.DeviceType,
					Host:       asaDevice.Host,
					Port:       asaDevice.Port,
					Ipv4:       asaDevice.Ipv4,
					SdcType:    asaDevice.LarType,
					SdcUid:     asaDevice.LarUid,
				}
				assert.Equal(t, expectedCreatedOutput, *output)

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDevice.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "successfully onboards ASA when using CDG after recovering from certificate error",

			input: asa.CreateInput{
				Name:             asaDevice.Name,
				SdcType:          asaDevice.LarType,
				Ipv4:             asaDevice.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDevice)
				configureDeviceReadSpecificToRespondSuccessfully(asaDevice.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: asaconfig.AsaConfigStateError,
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
					Uid:        asaDevice.Uid,
					Name:       asaDevice.Name,
					DeviceType: asaDevice.DeviceType,
					Host:       asaDevice.Host,
					Port:       asaDevice.Port,
					Ipv4:       asaDevice.Ipv4,
					SdcType:    asaDevice.LarType,
					SdcUid:     asaDevice.LarUid,
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
			testName: "successfully onboards ASA when using SDC",
			input: asa.CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				SdcType:          asaDeviceUsingSdc.LarType,
				SdcUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: false,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureSdcReadToRespondSuccessfully(sdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedCreatedOutput := asa.CreateOutput{
					Uid:        asaDeviceUsingSdc.Uid,
					Name:       asaDeviceUsingSdc.Name,
					DeviceType: asaDeviceUsingSdc.DeviceType,
					Host:       asaDeviceUsingSdc.Host,
					Port:       asaDeviceUsingSdc.Port,
					Ipv4:       asaDeviceUsingSdc.Ipv4,
					SdcType:    asaDeviceUsingSdc.LarType,
					SdcUid:     asaDeviceUsingSdc.LarUid,
				}
				assert.Equal(t, expectedCreatedOutput, *output)

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDeviceUsingSdc.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertSdcReadByUidWasCalledOnce(sdc.Uid, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "successfully onboards ASA when using SDC after recovering from certificate error",

			input: asa.CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				SdcType:          asaDeviceUsingSdc.LarType,
				SdcUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: asaconfig.AsaConfigStateError,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureSdcReadToRespondSuccessfully(sdc)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedCreatedOutput := asa.CreateOutput{
					Uid:        asaDeviceUsingSdc.Uid,
					Name:       asaDeviceUsingSdc.Name,
					DeviceType: asaDeviceUsingSdc.DeviceType,
					Host:       asaDeviceUsingSdc.Host,
					Port:       asaDeviceUsingSdc.Port,
					Ipv4:       asaDeviceUsingSdc.Ipv4,
					SdcType:    asaDeviceUsingSdc.LarType,
					SdcUid:     asaDeviceUsingSdc.LarUid,
				}
				assert.Equal(t, expectedCreatedOutput, *output)

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDeviceUsingSdc.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertDeviceUpdateWasCalledOnce(asaDeviceUsingSdc.Uid, t)
				assertSdcReadByUidWasCalledOnce(sdc.Uid, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "returns error when onboarding ASA and create device call experiences issues",

			input: asa.CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				SdcType:          asaDeviceUsingSdc.LarType,
				SdcUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondWithError()
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: asaconfig.AsaConfigStateError,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureSdcReadToRespondSuccessfully(sdc)
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
				Name:             asaDeviceUsingSdc.Name,
				SdcType:          asaDeviceUsingSdc.LarType,
				SdcUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondWithError(asaDeviceUsingSdc.Uid)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: asaconfig.AsaConfigStateError,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureSdcReadToRespondSuccessfully(sdc)
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
				Name:             asaDeviceUsingSdc.Name,
				SdcType:          asaDeviceUsingSdc.LarType,
				SdcUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithError(asaConfig.Uid)
				configureSdcReadToRespondSuccessfully(sdc)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *asa.CreateOutput, err *asa.CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when onboarding ASA and sdc read call experiences issues",

			input: asa.CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				SdcType:          asaDeviceUsingSdc.LarType,
				SdcUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: asaconfig.AsaConfigStateError,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureSdcReadToRespondWithError(sdc.Uid)
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
				Name:             asaDeviceUsingSdc.Name,
				SdcType:          asaDeviceUsingSdc.LarType,
				SdcUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: asaconfig.AsaConfigStateError,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureSdcReadToRespondSuccessfully(sdc)
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
				Name:             asaDeviceUsingSdc.Name,
				SdcType:          asaDeviceUsingSdc.LarType,
				SdcUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input asa.CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithCalls(asaConfig.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, asaconfig.ReadOutput{
						Uid:   asaSpecificDevice.SpecificUid,
						State: asaconfig.AsaConfigStateError,
					}),
					httpmock.NewJsonResponderOrPanic(200, asaConfig),
				})
				configureSdcReadToRespondSuccessfully(sdc)
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

			output, err := asa.Create(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
