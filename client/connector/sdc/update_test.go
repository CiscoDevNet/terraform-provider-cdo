package sdc_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validUpdateRequestOutput := sdc.NewUpdateResponseBuilder().
		Uid(sdcUid).
		Name(sdcName).
		Build()

	validUpdateOutput := sdc.NewUpdateOutputBuilder().
		UpdateRequestOutput(validUpdateRequestOutput).
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
		sdcUid     string
		sdcName    string
		setupFunc  func()
		assertFunc func(output *sdc.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully update SDC",
			sdcUid:   sdcUid,
			sdcName:  sdcName,

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", sdcUid),
					httpmock.NewJsonResponderOrPanic(200, validUpdateOutput),
				)
				httpmock.RegisterResponder(
					http.MethodPost,
					fmt.Sprintf("/anubis/rest/v1/oauth/token"),
					httpmock.NewJsonResponderOrPanic(200, validUserToken),
				)
			},

			assertFunc: func(output *sdc.UpdateOutput, err error, t *testing.T) {
				internalTesting.AssertNil(t, err, "error should be nil")
				internalTesting.AssertNotNil(t, output, "output should not be nil")
				internalTesting.AssertDeepEqual(t, validUpdateOutput, *output, "output should be same as valid output")
			},
		},
		{
			testName: "should error if failed to update sdc",
			sdcUid:   sdcUid,
			sdcName:  sdcName,

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", sdcUid),
					httpmock.NewJsonResponderOrPanic(500, "test error"),
				)
				httpmock.RegisterResponder(
					http.MethodPost,
					fmt.Sprintf("/anubis/rest/v1/oauth/token"),
					httpmock.NewJsonResponderOrPanic(200, validUserToken),
				)
			},

			assertFunc: func(output *sdc.UpdateOutput, err error, t *testing.T) {
				internalTesting.AssertNotNil(t, err, "error should not be nil")
				internalTesting.AssertDeepEqual(t, output, &sdc.UpdateOutput{}, "output should be zero value")
			},
		},
		{
			testName: "should error if failed to generate bootstrap data",
			sdcUid:   sdcUid,
			sdcName:  sdcName,

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", sdcUid),
					httpmock.NewJsonResponderOrPanic(200, validUpdateOutput),
				)
				httpmock.RegisterResponder(
					http.MethodPost,
					fmt.Sprintf("/anubis/rest/v1/oauth/token"),
					httpmock.NewJsonResponderOrPanic(500, nil),
				)
			},

			assertFunc: func(output *sdc.UpdateOutput, err error, t *testing.T) {
				internalTesting.AssertNotNil(t, err, "error should not be nil")
				internalTesting.AssertDeepEqual(t, output, &sdc.UpdateOutput{}, "output should be zero value")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := sdc.Update(context.Background(), *internalHttp.MustNewWithDefault(baseUrl, "a_valid_token"), sdc.NewUpdateInput(testCase.sdcUid, testCase.sdcName))

			testCase.assertFunc(output, err, t)
		})
	}
}
