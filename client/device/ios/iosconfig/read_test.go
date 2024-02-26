package iosconfig_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios/iosconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestIosConfigRead(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validIosConfig := iosconfig.ReadOutput{
		Uid:   iosConfigUid,
		State: state.DONE,
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *iosconfig.ReadOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully fetch iOS config",
			targetUid: iosConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					buildIosConfigPath(iosConfigUid),
					httpmock.NewJsonResponderOrPanic(200, validIosConfig),
				)
			},

			assertFunc: func(output *iosconfig.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validIosConfig, *output)
			},
		},
		{
			testName:  "returns nil output when iOS config not found",
			targetUid: iosConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					buildIosConfigPath(iosConfigUid),
					httpmock.NewStringResponder(404, ""),
				)
			},

			assertFunc: func(output *iosconfig.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
		{
			testName:  "return error when fetching iOS Config and remote service encounters issue",
			targetUid: iosConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					buildIosConfigPath(iosConfigUid),
					httpmock.NewStringResponder(500, "service is experiencing issues"),
				)
			},

			assertFunc: func(output *iosconfig.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := iosconfig.Read(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				*iosconfig.NewReadInput(iosConfigUid),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
