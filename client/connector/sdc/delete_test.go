package sdc_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/http"
	internalTesting "github.com/CiscoDevNet/terraform-provider-cdo/go-client/internal/testing"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validDeleteOutput := sdc.DeleteOutput{}

	testCases := []struct {
		testName   string
		sdcUid     string
		setupFunc  func()
		assertFunc func(output *sdc.DeleteOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully delete SDC",
			sdcUid:   sdcUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"DELETE",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", sdcUid),
					httpmock.NewJsonResponderOrPanic(200, validDeleteOutput),
				)
			},

			assertFunc: func(output *sdc.DeleteOutput, err error, t *testing.T) {
				internalTesting.AssertNil(t, err, "error should be nil")
				internalTesting.AssertNotNil(t, output, "output should not be nil")
				internalTesting.AssertDeepEqual(t, validDeleteOutput, *output, "output should be valid")
			},
		},
		{
			testName: "should error if failed to delete",
			sdcUid:   sdcUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"DELETE",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", sdcUid),
					httpmock.NewJsonResponderOrPanic(500, nil),
				)
			},

			assertFunc: func(output *sdc.DeleteOutput, err error, t *testing.T) {
				internalTesting.AssertNotNil(t, err, "error should be nil")
				internalTesting.AssertDeepEqual(t, &sdc.DeleteOutput{}, output, "output should be zero value")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := sdc.Delete(context.Background(), *http.NewWithDefault(baseUrl, "a_valid_token"), sdc.NewDeleteInput(testCase.sdcUid))

			testCase.assertFunc(output, err, t)
		})
	}
}
