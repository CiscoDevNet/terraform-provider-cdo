package sdc

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestReadByUid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validCdg := NewSdcResponseBuilder().
		AsDefaultCloudConnector().
		WithUid(cdgUid).
		WithName(cdgName).
		WithTenantUid(tenantUid).
		Build()

	validSdc := NewSdcResponseBuilder().
		AsOnPremConnector().
		WithUid(sdcUid).
		WithName(sdcName).
		WithTenantUid(tenantUid).
		Build()

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *ReadOutput, err error, t *testing.T)
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

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
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

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
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

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got (dereferenced): %+v", *output)
				}

				if err != nil {
					t.Errorf("expected err to be nil, got: %s", err.Error())
				}
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

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
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

			output, err := ReadByUid(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), ReadInput{LarUid: testCase.targetUid})

			testCase.assertFunc(output, err, t)
		})
	}
}

func TestReadByName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validCdg := NewSdcResponseBuilder().
		AsDefaultCloudConnector().
		WithUid(cdgUid).
		WithName(cdgName).
		WithTenantUid(tenantUid).
		Build()

	validSdc := NewSdcResponseBuilder().
		AsOnPremConnector().
		WithUid(sdcUid).
		WithName(sdcName).
		WithTenantUid(tenantUid).
		Build()

	testCases := []struct {
		testName   string
		targetName string
		setupFunc  func()
		assertFunc func(output *ReadOutput, err error, t *testing.T)
	}{
		{
			testName:   "successfully fetch CDG by name",
			targetName: cdgName,

			setupFunc: func() {
				httpmock.RegisterResponderWithQuery(
					"GET",
					"/aegis/rest/v1/services/targets/proxies",
					fmt.Sprintf("q=name:%s", cdgName),
					httpmock.NewJsonResponderOrPanic(200, []ReadOutput{validCdg}),
				)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
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
					httpmock.NewJsonResponderOrPanic(200, []ReadOutput{validSdc}),
				)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
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
					httpmock.NewJsonResponderOrPanic(200, []ReadOutput{validCdg, validCdg}),
				)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
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
					httpmock.NewJsonResponderOrPanic(200, []ReadOutput{}),
				)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got (dereferenced): %+v", *output)
				}

				if err != nil {
					t.Errorf("expected err to be nil, got: %s", err.Error())
				}
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

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
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

			output, err := ReadByName(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), ReadByNameInput{LarName: testCase.targetName})

			testCase.assertFunc(output, err, t)
		})
	}
}
