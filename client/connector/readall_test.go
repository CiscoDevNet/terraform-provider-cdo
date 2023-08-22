package connector_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestReadAll(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validCdg := connector.NewConnectorOutputBuilder().
		AsDefaultCloudConnector().
		WithUid(cdgUid).
		WithName(cdgName).
		WithTenantUid(tenantUid).
		Build()

	validConnector := connector.NewConnectorOutputBuilder().
		AsOnPremConnector().
		WithUid(connectorUid).
		WithName(connectorName).
		WithTenantUid(tenantUid).
		Build()

	testCases := []struct {
		testName   string
		setupFunc  func()
		assertFunc func(output *connector.ReadAllOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully fetches secure connectors",

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					httpmock.NewJsonResponderOrPanic(200, connector.ReadAllOutput{validCdg, validConnector}),
				)
			},

			assertFunc: func(output *connector.ReadAllOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedResponse := connector.ReadAllOutput{validCdg, validConnector}
				assert.Equal(t, expectedResponse, *output)
			},
		},
		{
			testName: "returns empty slice no connectors have been onboarded",

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					httpmock.NewJsonResponderOrPanic(200, connector.ReadAllOutput{}),
				)
			},

			assertFunc: func(output *connector.ReadAllOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Len(t, *output, 0)
			},
		},
		{
			testName: "return error when fetching all secure connectors and remote service encounters issue",

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					httpmock.NewStringResponder(500, "service is experiencing issues"),
				)
			},

			assertFunc: func(output *connector.ReadAllOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := connector.ReadAll(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				*connector.NewReadAllInput(),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
