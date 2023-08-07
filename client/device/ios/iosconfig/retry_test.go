package iosconfig

import (
	"context"
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
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				internalTesting.AssertEndpointCalledTimes("GET", buildIosConfigPath(iosConfigUid), 3, t)
			},
		},

		{
			testName:  "returns error if config state transitions to error",
			targetUid: iosConfigUid,

			setupFunc: func() {
				errorIosConfig := ReadOutput{
					Uid:   iosConfigUid,
					State: IosConfigStateError,
				}

				configureIosConfigReadToSucceedInSubsequentCalls(iosConfigUid, []ReadOutput{
					inProgressIosConfig,
					inProgressIosConfig,
					errorIosConfig,
				})
			},

			assertFunc: func(err error, t *testing.T) {
				if err == nil {
					t.Error("expected error to be returned")
				}

				internalTesting.AssertEndpointCalledTimes("GET", buildIosConfigPath(iosConfigUid), 3, t)
			},
		},

		{
			testName:  "return errors if config state transitions to bad credentials",
			targetUid: iosConfigUid,

			setupFunc: func() {
				badCredentialsIosConfig := ReadOutput{
					Uid:   iosConfigUid,
					State: IosConfigStateBadCredentials,
				}

				configureIosConfigReadToSucceedInSubsequentCalls(iosConfigUid, []ReadOutput{
					inProgressIosConfig,
					inProgressIosConfig,
					badCredentialsIosConfig,
				})
			},

			assertFunc: func(err error, t *testing.T) {
				if err == nil {
					t.Error("expected error to be returned")
				}

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

			err := retry.Do(UntilState(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.targetUid, IosConfigStateDone), retryOptions)

			testCase.assertFunc(err, t)
		})
	}
}
