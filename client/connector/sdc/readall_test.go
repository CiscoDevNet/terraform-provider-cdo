package sdc_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestReadAll(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validCdg := sdc.NewSdcOutputBuilder().
		AsDefaultCloudConnector().
		WithUid(cdgUid).
		WithName(cdgName).
		WithTenantUid(tenantUid).
		Build()

	validSdc := sdc.NewSdcOutputBuilder().
		AsOnPremConnector().
		WithUid(sdcUid).
		WithName(sdcName).
		WithTenantUid(tenantUid).
		Build()

	testCases := []struct {
		testName   string
		setupFunc  func()
		assertFunc func(output *sdc.ReadAllOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully fetches secure connectors",

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					httpmock.NewJsonResponderOrPanic(200, sdc.ReadAllOutput{validCdg, validSdc}),
				)
			},

			assertFunc: func(output *sdc.ReadAllOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)

				expectedResponse := sdc.ReadAllOutput{validCdg, validSdc}
				assert.Equal(t, expectedResponse, *output)
			},
		},
		{
			testName: "returns empty slice no connectors have been onboarded",

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					httpmock.NewJsonResponderOrPanic(200, sdc.ReadAllOutput{}),
				)
			},

			assertFunc: func(output *sdc.ReadAllOutput, err error, t *testing.T) {
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

			assertFunc: func(output *sdc.ReadAllOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := sdc.ReadAll(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				*sdc.NewReadAllInput(),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
