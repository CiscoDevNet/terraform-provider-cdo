package sdc_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestReadByUid(t *testing.T) {
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
		targetUid  string
		setupFunc  func()
		assertFunc func(output *sdc.ReadOutput, err error, t *testing.T)
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

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(validCdg, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validCdg, output)
				}
			},
		},
		{
			testName:  "successfully fetch SDC by uid",
			targetUid: sdcUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", sdcUid),
					httpmock.NewJsonResponderOrPanic(200, validSdc),
				)
			},

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(validSdc, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validSdc, output)
				}
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

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got (dereferenced): %+v", *output)
				}

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

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got (dereferenced): %+v", *output)
				}

				if err == nil {
					t.Error("error was nil!")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := sdc.ReadByUid(context.Background(), *http.MustNewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), sdc.ReadByUidInput{SdcUid: testCase.targetUid})

			testCase.assertFunc(output, err, t)
		})
	}
}

func TestReadByName(t *testing.T) {
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
		targetName string
		setupFunc  func()
		assertFunc func(output *sdc.ReadOutput, err error, t *testing.T)
	}{
		{
			testName:   "successfully fetch CDG by name",
			targetName: cdgName,

			setupFunc: func() {
				httpmock.RegisterResponderWithQuery(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					fmt.Sprintf("q=name:%s", cdgName),
					httpmock.NewJsonResponderOrPanic(200, []sdc.ReadOutput{validCdg}),
				)
			},

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(validCdg, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validCdg, output)
				}
			},
		},
		{
			testName:   "successfully fetch SDC by name",
			targetName: sdcName,

			setupFunc: func() {
				httpmock.RegisterResponderWithQuery(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					fmt.Sprintf("q=name:%s", sdcName),
					httpmock.NewJsonResponderOrPanic(200, []sdc.ReadOutput{validSdc}),
				)
			},

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(validSdc, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validSdc, output)
				}
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
					httpmock.NewJsonResponderOrPanic(200, []sdc.ReadOutput{validCdg, validCdg}),
				)
			},

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", output)
				}

				if err == nil {
					t.Errorf("expected error")
				}

				expectedError := fmt.Sprintf("multiple SDCs found with the name: %s", cdgName)
				if err.Error() != expectedError {
					t.Errorf("expected error: '%s', got: '%s'", expectedError, err.Error())
				}
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
					httpmock.NewJsonResponderOrPanic(200, []sdc.ReadOutput{}),
				)
			},

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got (dereferenced): %+v", *output)
				}

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

			assertFunc: func(output *sdc.ReadOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got (dereferenced): %+v", *output)
				}

				if err == nil {
					t.Error("error was nil!")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := sdc.ReadByName(context.Background(), *http.MustNewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), sdc.ReadByNameInput{SdcName: testCase.targetName})

			testCase.assertFunc(output, err, t)
		})
	}
}
