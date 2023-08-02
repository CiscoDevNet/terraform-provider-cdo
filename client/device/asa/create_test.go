package asa

import (
	"context"
	"reflect"
	"testing"

	"github.com/cisco-lockhart/go-client/connector/sdc"
	"github.com/cisco-lockhart/go-client/device"
	"github.com/cisco-lockhart/go-client/internal/device/asaconfig"
	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestAsaCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	asaDevice := device.NewReadOutputBuilder().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-asa").
		OnboardedUsingCdg("88888888-8888-8888-8888-888888888888").
		WithLocation("10.10.0.1", 443).
		Build()

	asaDeviceUsingSdc := device.NewReadOutputBuilder().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-asa").
		OnboardedUsingSdc("99999999-9999-9999-9999-999999999999").
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

	sdc := sdc.NewSdcResponseBuilder().
		WithName("CloudDeviceGateway").
		WithUid(asaDeviceUsingSdc.LarUid).
		WithTenantUid("44444444-4444-4444-4444-444444444444").
		AsOnPremConnector().
		Build()

	testCases := []struct {
		testName   string
		input      CreateInput
		setupFunc  func(input CreateInput)
		assertFunc func(output *CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully onboards ASA when using CDG",
			input: CreateInput{
				Name:             asaDevice.Name,
				LarType:          asaDevice.LarType,
				Ipv4:             asaDevice.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: false,
			},

			setupFunc: func(input CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDevice)
				configureDeviceReadSpecificToRespondSuccessfully(asaDevice.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatalf("output is nil!")
				}

				expectedCreatedOutput := CreateOutput{
					Uid:         asaDevice.Uid,
					Name:        asaDevice.Name,
					DeviceType:  asaDevice.DeviceType,
					Host:        asaDevice.Host,
					Port:        asaDevice.Port,
					Ipv4:        asaDevice.Ipv4,
					LarType:     asaDevice.LarType,
					LarUid:      asaDevice.LarUid,
					specificUid: asaConfig.Uid,
				}
				if !reflect.DeepEqual(expectedCreatedOutput, *output) {
					t.Errorf("expected: %+v, got: %+v", expectedCreatedOutput, output)
				}

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDevice.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "successfully onboards ASA when using CDG after recovering from certificate error",

			input: CreateInput{
				Name:             asaDevice.Name,
				LarType:          asaDevice.LarType,
				Ipv4:             asaDevice.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input CreateInput) {
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

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatalf("output is nil!")
				}

				expectedCreatedOutput := CreateOutput{
					Uid:         asaDevice.Uid,
					Name:        asaDevice.Name,
					DeviceType:  asaDevice.DeviceType,
					Host:        asaDevice.Host,
					Port:        asaDevice.Port,
					Ipv4:        asaDevice.Ipv4,
					LarType:     asaDevice.LarType,
					LarUid:      asaDevice.LarUid,
					specificUid: asaConfig.Uid,
				}
				if !reflect.DeepEqual(expectedCreatedOutput, *output) {
					t.Errorf("expected: %+v, got: %+v", expectedCreatedOutput, output)
				}

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDevice.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertDeviceUpdateWasCalledOnce(asaDevice.Uid, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "successfully onboards ASA when using SDC",
			input: CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				LarType:          asaDeviceUsingSdc.LarType,
				LarUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: false,
			},

			setupFunc: func(input CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondSuccessfully(asaSpecificDevice.SpecificUid, asaConfig)
				configureSdcReadToRespondSuccessfully(sdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatalf("output is nil!")
				}

				expectedCreatedOutput := CreateOutput{
					Uid:         asaDeviceUsingSdc.Uid,
					Name:        asaDeviceUsingSdc.Name,
					DeviceType:  asaDeviceUsingSdc.DeviceType,
					Host:        asaDeviceUsingSdc.Host,
					Port:        asaDeviceUsingSdc.Port,
					Ipv4:        asaDeviceUsingSdc.Ipv4,
					LarType:     asaDeviceUsingSdc.LarType,
					LarUid:      asaDeviceUsingSdc.LarUid,
					specificUid: asaConfig.Uid,
				}
				if !reflect.DeepEqual(expectedCreatedOutput, *output) {
					t.Errorf("expected: %+v, got: %+v", expectedCreatedOutput, output)
				}

				assertDeviceCreateWasCalledOnce(t)
				assertDeviceReadSpecificWasCalledOnce(asaDeviceUsingSdc.Uid, t)
				assertAsaConfigReadWasCalledTimes(asaConfig.Uid, 2, t)
				assertSdcReadByUidWasCalledOnce(sdc.Uid, t)
				assertAsaConfigUpdateWasCalledOnce(asaConfig.Uid, t)
			},
		},

		{
			testName: "successfully onboards ASA when using SDC after recovering from certificate error",

			input: CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				LarType:          asaDeviceUsingSdc.LarType,
				LarUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input CreateInput) {
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

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatalf("output is nil!")
				}

				expectedCreatedOutput := CreateOutput{
					Uid:         asaDeviceUsingSdc.Uid,
					Name:        asaDeviceUsingSdc.Name,
					DeviceType:  asaDeviceUsingSdc.DeviceType,
					Host:        asaDeviceUsingSdc.Host,
					Port:        asaDeviceUsingSdc.Port,
					Ipv4:        asaDeviceUsingSdc.Ipv4,
					LarType:     asaDeviceUsingSdc.LarType,
					LarUid:      asaDeviceUsingSdc.LarUid,
					specificUid: asaConfig.Uid,
				}
				if !reflect.DeepEqual(expectedCreatedOutput, *output) {
					t.Errorf("expected: %+v, got: %+v", expectedCreatedOutput, output)
				}

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

			input: CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				LarType:          asaDeviceUsingSdc.LarType,
				LarUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input CreateInput) {
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

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err == nil {
					t.Error("error is nil!")
				}

				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", output)
				}
			},
		},

		{
			testName: "returns error when onboarding ASA and read specific device call experiences issues",

			input: CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				LarType:          asaDeviceUsingSdc.LarType,
				LarUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input CreateInput) {
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

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err == nil {
					t.Error("error is nil!")
				}

				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", output)
				}
			},
		},

		{
			testName: "returns error when onboarding ASA and read asa config call experiences issues",

			input: CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				LarType:          asaDeviceUsingSdc.LarType,
				LarUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input CreateInput) {
				configureDeviceCreateToRespondSuccessfully(asaDeviceUsingSdc)
				configureDeviceReadSpecificToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaSpecificDevice)
				configureAsaConfigReadToRespondWithError(asaConfig.Uid)
				configureSdcReadToRespondSuccessfully(sdc)
				configureDeviceUpdateToRespondSuccessfully(asaDeviceUsingSdc.Uid, asaDeviceUsingSdc)
				configureAsaConfigUpdateToRespondSuccessfully(asaConfig.Uid, asaconfig.UpdateOutput{Uid: asaConfig.Uid})
			},

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err == nil {
					t.Error("error is nil!")
				}

				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", output)
				}
			},
		},

		{
			testName: "returns error when onboarding ASA and sdc read call experiences issues",

			input: CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				LarType:          asaDeviceUsingSdc.LarType,
				LarUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input CreateInput) {
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

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err == nil {
					t.Error("error is nil!")
				}

				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", output)
				}
			},
		},

		{
			testName: "returns error when onboarding ASA and device update call experiences issues",

			input: CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				LarType:          asaDeviceUsingSdc.LarType,
				LarUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input CreateInput) {
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

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
				if err == nil {
					t.Error("error is nil!")
				}

				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", output)
				}
			},
		},

		{
			testName: "returns error when onboarding ASA and ASA config update call experiences issues",

			input: CreateInput{
				Name:             asaDeviceUsingSdc.Name,
				LarType:          asaDeviceUsingSdc.LarType,
				LarUid:           asaDeviceUsingSdc.LarUid,
				Ipv4:             asaDeviceUsingSdc.Ipv4,
				Username:         "unittestuser",
				Password:         "not a real password",
				IgnoreCertifcate: true,
			},

			setupFunc: func(input CreateInput) {
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

			assertFunc: func(output *CreateOutput, err error, t *testing.T) {
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

			testCase.setupFunc(testCase.input)

			output, err := Create(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
