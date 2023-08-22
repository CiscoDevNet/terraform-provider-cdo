package connector_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validDeleteOutput := connector.DeleteOutput{}

	testCases := []struct {
		testName   string
		sdcUid     string
		setupFunc  func()
		assertFunc func(output *connector.DeleteOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully delete SDC",
			sdcUid:   connectorUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"DELETE",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", connectorUid),
					httpmock.NewJsonResponderOrPanic(200, validDeleteOutput),
				)
			},

			assertFunc: func(output *connector.DeleteOutput, err error, t *testing.T) {
				assert.Nil(t, err, "error should be nil")
				assert.NotNil(t, output, "output should not be nil")
				assert.Equal(t, validDeleteOutput, *output, "output should be valid")
			},
		},
		{
			testName: "should error if failed to delete",
			sdcUid:   connectorUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"DELETE",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", connectorUid),
					httpmock.NewJsonResponderOrPanic(500, nil),
				)
			},

			assertFunc: func(output *connector.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err, "error should be nil")
				assert.Equal(t, &connector.DeleteOutput{}, output, "output should be zero value")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := connector.Delete(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				connector.NewDeleteInput(testCase.sdcUid),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
