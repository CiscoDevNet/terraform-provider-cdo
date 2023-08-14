package iosconfig

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/jarcoal/httpmock"
)

func TestIosConfigUntilState(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validIosConfig := ReadOutput{
		Uid:   iosConfigUid,
		State: IosConfigStateDone,
	}

	inProgressIosConfig := ReadOutput{
		Uid:   iosConfigUid,
		State: "SOME_INTERMEDIATE_STATE",
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(err error, t *testing.T)
	}{
		{
			testName:  "successfully returns once state reaches done",
			targetUid: iosConfigUid,

			setupFunc: func() {
				configureIosConfigReadToSucceedInSubsequentCalls(iosConfigUid, []ReadOutput{
					inProgressIosConfig,
					inProgressIosConfig,
					validIosConfig,
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)
				internalTesting.AssertEndpointCalledTimes("GET", buildIosConfigPath(iosConfigUid), 3, t)
			},
		},

		{
			testName:  "returns error if config state transitions to error",
			targetUid: iosConfigUid,

			setupFunc: func() {
				errorIosConfig := ReadOutput{
					Uid:   iosConfigUid,
					State: state.ERROR,
				}

				configureIosConfigReadToSucceedInSubsequentCalls(iosConfigUid, []ReadOutput{
					inProgressIosConfig,
					inProgressIosConfig,
					errorIosConfig,
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
				internalTesting.AssertEndpointCalledTimes("GET", buildIosConfigPath(iosConfigUid), 3, t)
			},
		},

		{
			testName:  "return errors if config state transitions to bad credentials",
			targetUid: iosConfigUid,

			setupFunc: func() {
				badCredentialsIosConfig := ReadOutput{
					Uid:   iosConfigUid,
					State: state.BAD_CREDENTIALS,
				}

				configureIosConfigReadToSucceedInSubsequentCalls(iosConfigUid, []ReadOutput{
					inProgressIosConfig,
					inProgressIosConfig,
					badCredentialsIosConfig,
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
				internalTesting.AssertEndpointCalledTimes("GET", buildIosConfigPath(iosConfigUid), 3, t)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			retryOptions := retry.DefaultOpts
			retryOptions.Delay = 1 * time.Millisecond

			err := retry.Do(UntilState(context.Background(), *http.MustNewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.targetUid, state.ERROR), retryOptions)

			testCase.assertFunc(err, t)
		})
	}
}
