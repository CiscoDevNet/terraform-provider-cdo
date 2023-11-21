package connector_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validDeleteOutput := connector.DeleteOutput{}
	validConnector := connector.NewConnectorOutputBuilder().
		WithUid(connectorUid).
		WithTenantUid(tenantUid).
		WithName(connectorName).
		Build()
	connectorPresentResponse := httpmock.NewJsonResponderOrPanic(200, []connector.ReadOutput{validConnector})
	connectorMissingResponse := httpmock.NewJsonResponderOrPanic(200, []connector.ReadOutput{})
	connectorDeleteInitiatedResponse := httpmock.NewJsonResponderOrPanic(200, validDeleteOutput)
	errorMessage := "intentional failed to delete error"

	testCases := []struct {
		testName   string
		sdcUid     string
		setupFunc  func()
		assertFunc func(output *connector.DeleteOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully delete SDC",
			sdcUid:   validConnector.Uid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAllConnectors(baseUrl),
					connectorPresentResponse.Then(connectorMissingResponse),
				)
				httpmock.RegisterResponder(
					http.MethodDelete,
					url.DeleteConnector(baseUrl, validConnector.Uid),
					connectorDeleteInitiatedResponse,
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
			sdcUid:   validConnector.Uid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadAllConnectors(baseUrl),
					connectorPresentResponse.Then(connectorMissingResponse),
				)
				httpmock.RegisterResponder(
					http.MethodDelete,
					url.DeleteConnector(baseUrl, validConnector.Uid),
					httpmock.NewJsonResponderOrPanic(500, errorMessage),
				)
			},

			assertFunc: func(output *connector.DeleteOutput, err error, t *testing.T) {
				assert.NotNil(t, err, "error should be nil")
				assert.ErrorContains(t, err, errorMessage)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := connector.Delete(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				connector.NewDeleteInput(testCase.sdcUid),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
