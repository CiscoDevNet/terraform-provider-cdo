package connector_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestReadByUid(t *testing.T) {
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
		targetUid  string
		setupFunc  func()
		assertFunc func(output *connector.ReadOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully fetch CDG by uid",
			targetUid: cdgUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", cdgUid),
					httpmock.NewJsonResponderOrPanic(200, validCdg),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validCdg, *output)
			},
		},
		{
			testName:  "successfully fetch connector by uid",
			targetUid: connectorUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", connectorUid),
					httpmock.NewJsonResponderOrPanic(200, validConnector),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validConnector, *output)
			},
		},
		{
			testName:  "returns nil ouput when CDG not found",
			targetUid: cdgUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", cdgUid),
					httpmock.NewStringResponder(404, ""),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
		{
			testName:  "return error when fetching CDG and remote service encounters issue",
			targetUid: cdgUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", cdgUid),
					httpmock.NewStringResponder(500, "service is experiencing issues"),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := connector.ReadByUid(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				*connector.NewReadByUidInput(testCase.targetUid),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func TestReadByName(t *testing.T) {
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
		targetName string
		setupFunc  func()
		assertFunc func(output *connector.ReadOutput, err error, t *testing.T)
	}{
		{
			testName:   "successfully fetch CDG by name",
			targetName: cdgName,

			setupFunc: func() {
				httpmock.RegisterResponderWithQuery(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					fmt.Sprintf("q=name:%s", cdgName),
					httpmock.NewJsonResponderOrPanic(200, []connector.ReadOutput{validCdg}),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validCdg, *output)
			},
		},
		{
			testName:   "successfully fetch connector by name",
			targetName: connectorName,

			setupFunc: func() {
				httpmock.RegisterResponderWithQuery(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					fmt.Sprintf("q=name:%s", connectorName),
					httpmock.NewJsonResponderOrPanic(200, []connector.ReadOutput{validConnector}),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validConnector, *output)
			},
		},
		{
			testName:   "returns error when response from remote service returns multiple CDGs",
			targetName: cdgName,

			setupFunc: func() {
				httpmock.RegisterResponderWithQuery(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					fmt.Sprintf("q=name:%s", cdgName),
					httpmock.NewJsonResponderOrPanic(200, []connector.ReadOutput{validCdg, validCdg}),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)

				expectedError := fmt.Sprintf("multiple connector found with the name: %s", cdgName)
				assert.Equal(t, expectedError, err.Error())
			},
		},
		{
			testName:   "returns nil ouput when CDG not found",
			targetName: cdgName,

			setupFunc: func() {
				httpmock.RegisterResponderWithQuery(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					fmt.Sprintf("q=name:%s", cdgName),
					httpmock.NewJsonResponderOrPanic(200, []connector.ReadOutput{}),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
		{
			testName:   "return error when fetching CDG and remote service encounters issue",
			targetName: cdgName,

			setupFunc: func() {
				httpmock.RegisterResponderWithQuery(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					fmt.Sprintf("q=name:%s", cdgName),
					httpmock.NewStringResponder(500, "service encountered issue"),
				)
			},

			assertFunc: func(output *connector.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := connector.ReadByName(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				*connector.NewReadByNameInput(testCase.targetName),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
