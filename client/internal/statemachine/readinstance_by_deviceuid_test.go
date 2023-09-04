package statemachine_test

import (
	"context"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestReadInstanceByDeviceUid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName   string
		input      statemachine.ReadInstanceByDeviceUidInput
		setupFunc  func()
		assertFunc func(output *statemachine.ReadInstanceByDeviceUidOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully read state machine instance by uid",
			input:    statemachine.NewReadInstanceByDeviceUidInput(deviceUid),
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadStateMachineInstance(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, []statemachine.ReadInstanceByDeviceUidOutput{
						validReadStateMachineOutput,
					}),
				)
			},
			assertFunc: func(output *statemachine.ReadInstanceByDeviceUidOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validReadStateMachineOutput, *output)
			},
		},
		{
			testName: "error when read state machine instance by uid error",
			input:    statemachine.NewReadInstanceByDeviceUidInput(deviceUid),
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadStateMachineInstance(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *statemachine.ReadInstanceByDeviceUidOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := statemachine.ReadInstanceByDeviceUid(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
