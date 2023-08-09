package iosconfig

import (
	"context"
	"reflect"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestIosConfigRead(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validIosConfig := ReadOutput{
		Uid:   iosConfigUid,
		State: IosConfigStateDone,
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *ReadOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully fetch iOS config",
			targetUid: iosConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					buildIosConfigPath(iosConfigUid),
					httpmock.NewJsonResponderOrPanic(200, validIosConfig),
				)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(validIosConfig, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validIosConfig, output)
				}
			},
		},
		{
			testName:  "returns nil ouput when iOS config not found",
			targetUid: iosConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					buildIosConfigPath(iosConfigUid),
					httpmock.NewStringResponder(404, ""),
				)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
		{
			testName:  "return error when fetching iOS Config and remote service encounters issue",
			targetUid: iosConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					buildIosConfigPath(iosConfigUid),
					httpmock.NewStringResponder(500, "service is experiencing issues"),
				)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got (dereferenced): %+v", *output)
				}

				if err == nil {
					t.Error("error is nil!")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := Read(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), *NewReadInput(iosConfigUid))

			testCase.assertFunc(output, err, t)
		})
	}
}
