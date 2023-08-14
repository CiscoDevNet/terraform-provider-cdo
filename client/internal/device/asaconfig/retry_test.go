package asaconfig

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	h "net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/jarcoal/httpmock"
)

func makeFetchAsaConfigUrl(asaconfigUid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/asa/configs/%s", asaConfigUid)
}

func TestAsaConfigUntilStateDone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validAsaConfig := ReadOutput{
		Uid:   asaConfigUid,
		State: AsaConfigStateDone,
	}

	inProgressAsaConfig := ReadOutput{
		Uid:   asaConfigUid,
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
			targetUid: asaConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					makeFetchAsaConfigUrl(asaConfigUid),
					httpmock.NewJsonResponderOrPanic(200, validAsaConfig),
				)

				callCount := 0
				httpmock.RegisterResponder("GET", makeFetchAsaConfigUrl(asaConfigUid), func(r *h.Request) (*h.Response, error) {
					callCount += 1

					if callCount < 3 {
						return httpmock.NewJsonResponse(200, inProgressAsaConfig)
					}

					return httpmock.NewJsonResponse(200, validAsaConfig)
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)

				internalTesting.AssertEndpointCalledTimes("GET", makeFetchAsaConfigUrl(asaConfigUid), 3, t)
			},
		},
		{
			testName:  "returns error if config state transitions to error",
			targetUid: asaConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					makeFetchAsaConfigUrl(asaConfigUid),
					httpmock.NewJsonResponderOrPanic(200, validAsaConfig),
				)

				errorAsaConfig := ReadOutput{
					Uid:   asaConfigUid,
					State: AsaConfigStateError,
				}

				callCount := 0
				httpmock.RegisterResponder("GET", makeFetchAsaConfigUrl(asaConfigUid), func(r *h.Request) (*h.Response, error) {
					callCount += 1

					if callCount < 3 {
						return httpmock.NewJsonResponse(200, inProgressAsaConfig)
					}

					return httpmock.NewJsonResponse(200, errorAsaConfig)
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)

				internalTesting.AssertEndpointCalledTimes("GET", makeFetchAsaConfigUrl(asaConfigUid), 3, t)
			},
		},
		{
			testName:  "return errors if config state transitions to bad credentials",
			targetUid: asaConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					makeFetchAsaConfigUrl(asaConfigUid),
					httpmock.NewJsonResponderOrPanic(200, validAsaConfig),
				)

				badCredentialsAsaConfig := ReadOutput{
					Uid:   asaConfigUid,
					State: AsaConfigStateBadCredentials,
				}

				callCount := 0
				httpmock.RegisterResponder("GET", makeFetchAsaConfigUrl(asaConfigUid), func(r *h.Request) (*h.Response, error) {
					callCount += 1

					if callCount < 3 {
						return httpmock.NewJsonResponse(200, inProgressAsaConfig)
					}

					return httpmock.NewJsonResponse(200, badCredentialsAsaConfig)
				})
			},

			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)

				internalTesting.AssertEndpointCalledTimes("GET", makeFetchAsaConfigUrl(asaConfigUid), 3, t)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			retryOptions := retry.DefaultOpts
			retryOptions.Delay = 1 * time.Millisecond

			err := retry.Do(UntilStateDone(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.targetUid), retryOptions)

			testCase.assertFunc(err, t)
		})
	}
}
