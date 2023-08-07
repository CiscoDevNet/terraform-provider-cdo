package sdc_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/http"
	internalTesting "github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/user"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validCreateRequestOutput := sdc.NewCreateResponseBuilder().
		Uid(sdcUid).
		TenantUid(tenantUid).
		Name(sdcName).
		ServiceConnectivityState(serviceConnectivityState).
		State(state).
		Status(status).
		Build()

	validCreateOutput := sdc.NewCreateOutputBuilder().
		CreateRequestOutput(validCreateRequestOutput).
		BootstrapData(bootstrapData).
		Build()

	validUserToken := user.NewGetTokenOutputBuilder().
		AccessToken(accessToken).
		RefreshToken(refreshToken).
		TenantUid(tenantUid).
		TenantName(tenantName).
		Scope(scope).
		TokenType(tokenType).
		Build()

	testCases := []struct {
		testName   string
		sdcName    string
		setupFunc  func()
		assertFunc func(output *sdc.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully create SDC",
			sdcName:  sdcName,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies"),
					httpmock.NewJsonResponderOrPanic(200, validCreateRequestOutput),
				)
				httpmock.RegisterResponder(
					"POST",
					fmt.Sprintf("/anubis/rest/v1/oauth/token"),
					httpmock.NewJsonResponderOrPanic(200, validUserToken),
				)
			},

			assertFunc: func(output *sdc.CreateOutput, err error, t *testing.T) {
				if err != nil {
					t.Fatalf("unexpected error: %s\n", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(validCreateOutput, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validCreateOutput, output)
				}
			},
		},
		{
			testName: "should error if failed to create proxy",
			sdcName:  sdcName,

			setupFunc: func() {
				internalTesting.MustResponseWithAtMostResponders(
					"POST",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies"),
					[]httpmock.Responder{
						httpmock.NewJsonResponderOrPanic(500, "test error"),
					},
				)
				httpmock.RegisterResponder(
					"POST",
					fmt.Sprintf("/anubis/rest/v1/oauth/token"),
					httpmock.NewJsonResponderOrPanic(200, validUserToken),
				)
			},

			assertFunc: func(output *sdc.CreateOutput, err error, t *testing.T) {
				internalTesting.AssertNotNil(t, err, "error should not be nil")
				internalTesting.AssertDeepEqual(t, output, &sdc.CreateOutput{}, "output should be zero value")
			},
		},
		{
			testName: "should error if failed to retrieve user token",
			sdcName:  sdcName,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies"),
					httpmock.NewJsonResponderOrPanic(200, validCreateRequestOutput),
				)
				internalTesting.MustResponseWithAtMostResponders(
					"POST",
					fmt.Sprintf("/anubis/rest/v1/oauth/token"),
					[]httpmock.Responder{
						httpmock.NewJsonResponderOrPanic(500, nil),
					},
				)
			},

			assertFunc: func(output *sdc.CreateOutput, err error, t *testing.T) {
				internalTesting.AssertNotNil(t, err, "error should not be nil")
				internalTesting.AssertDeepEqual(t, output, &sdc.CreateOutput{}, "output should be zero value")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := sdc.Create(context.Background(), *http.NewWithDefault(baseUrl, "a_valid_token"), *sdc.NewCreateInput(testCase.sdcName))

			testCase.assertFunc(output, err, t)
		})
	}
}