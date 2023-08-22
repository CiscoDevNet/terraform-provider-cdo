package ios

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios/iosconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestIosCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validConnector := connector.NewConnectorOutputBuilder().
		WithName("MyOnPremConnector").
		WithUid("88888888-8888-8888-8888-888888888888").
		WithTenantUid("00000000-0000-0000-0000-000000000000").
		AsOnPremConnector().
		Build()

	iosDevice := device.NewReadOutputBuilder().
		AsIos().
		WithUid("11111111-1111-1111-1111-111111111111").
		WithName("my-ios").
		OnboardedUsingOnPremConnector(validConnector.Uid).
		WithLocation("10.10.0.1", 443).
		Build()

	testCases := []struct {
		testName   string
		input      CreateInput
		setupFunc  func(input CreateInput)
		assertFunc func(output *CreateOutput, err *CreateError, t *testing.T)
	}{
		{
			testName: "successfully onboards iOS when using SDC",
			input: CreateInput{
				Name:              iosDevice.Name,
				ConnectorType:     iosDevice.ConnectorType,
				ConnectorUid:      iosDevice.ConnectorUid,
				SocketAddress:     iosDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input CreateInput) {
				configureDeviceCreateToRespondSuccessfully(iosDevice)
				configureSdcReadToRespondSuccessfully(validConnector)
				configureIosConfigReadToSucceedWithSubsequentCalls(iosDevice.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.PRE_READ_METADATA}),
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.DONE}),
				})
				configureDeviceUpdateToRespondSuccessfully(iosDevice)
			},

			assertFunc: func(output *CreateOutput, err *CreateError, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedCreatedOutput := CreateOutput{
					Uid:           iosDevice.Uid,
					Name:          iosDevice.Name,
					DeviceType:    iosDevice.DeviceType,
					Host:          iosDevice.Host,
					Port:          iosDevice.Port,
					SocketAddress: iosDevice.SocketAddress,
					ConnectorType: iosDevice.ConnectorType,
					ConnectorUid:  iosDevice.ConnectorUid,
				}
				assert.Equal(t, expectedCreatedOutput, *output)

				assertDeviceCreateWasCalledOnce(t)
				assertSdcReadByUidWasCalledOnce(validConnector.Uid, t)
				assertDeviceReadWasCalledTimes(iosDevice.Uid, 2, t)
				assertDeviceUpdateWasCalledOnce(iosDevice.Uid, t)
			},
		},

		{
			testName: "returns error when device create call encounters error",
			input: CreateInput{
				Name:              iosDevice.Name,
				ConnectorType:     iosDevice.ConnectorType,
				ConnectorUid:      iosDevice.ConnectorUid,
				SocketAddress:     iosDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input CreateInput) {
				configureDeviceCreateToRespondWithError()
				configureSdcReadToRespondSuccessfully(validConnector)
				configureIosConfigReadToSucceedWithSubsequentCalls(iosDevice.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.PRE_READ_METADATA}),
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.DONE}),
				})
				configureDeviceUpdateToRespondSuccessfully(iosDevice)
			},

			assertFunc: func(output *CreateOutput, err *CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when validConnector read call encounters error",
			input: CreateInput{
				Name:              iosDevice.Name,
				ConnectorType:     iosDevice.ConnectorType,
				ConnectorUid:      iosDevice.ConnectorUid,
				SocketAddress:     iosDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input CreateInput) {
				configureDeviceCreateToRespondSuccessfully(iosDevice)
				configureSdcReadToRespondWithError(validConnector.Uid)
				configureIosConfigReadToSucceedWithSubsequentCalls(iosDevice.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.PRE_READ_METADATA}),
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.DONE}),
				})
				configureDeviceUpdateToRespondSuccessfully(iosDevice)
			},

			assertFunc: func(output *CreateOutput, err *CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when iOS config read call encounters error",
			input: CreateInput{
				Name:              iosDevice.Name,
				ConnectorType:     iosDevice.ConnectorType,
				ConnectorUid:      iosDevice.ConnectorUid,
				SocketAddress:     iosDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input CreateInput) {
				configureDeviceCreateToRespondSuccessfully(iosDevice)
				configureSdcReadToRespondSuccessfully(validConnector)
				configureIosConfigReadToSucceedWithSubsequentCalls(iosDevice.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.ERROR}),
				})
				configureDeviceUpdateToRespondSuccessfully(iosDevice)
			},

			assertFunc: func(output *CreateOutput, err *CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},

		{
			testName: "returns error when device update call encounters error",
			input: CreateInput{
				Name:              iosDevice.Name,
				ConnectorType:     iosDevice.ConnectorType,
				ConnectorUid:      iosDevice.ConnectorUid,
				SocketAddress:     iosDevice.SocketAddress,
				Username:          "unittestuser",
				Password:          "not a real password",
				IgnoreCertificate: false,
			},

			setupFunc: func(input CreateInput) {
				configureDeviceCreateToRespondSuccessfully(iosDevice)
				configureSdcReadToRespondSuccessfully(validConnector)
				configureIosConfigReadToSucceedWithSubsequentCalls(iosDevice.Uid, []httpmock.Responder{
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.PRE_READ_METADATA}),
					httpmock.NewJsonResponderOrPanic(200, iosconfig.ReadOutput{Uid: iosDevice.Uid, State: state.DONE}),
				})
				configureDeviceUpdateToRespondWithError(iosDevice.Uid)
			},

			assertFunc: func(output *CreateOutput, err *CreateError, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := Create(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
