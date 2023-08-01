package asaconfig

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestAsaConfigReadByUid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	urlTemplate := "/aegis/rest/v1/services/asa/configs/%s"

	validAsaConfig := ReadOutput{
		Uid:   asaConfigUid,
		State: AsaConfigStateDone,
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *ReadOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully fetch ASA config",
			targetUid: asaConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf(urlTemplate, asaConfigUid),
					httpmock.NewJsonResponderOrPanic(200, validAsaConfig),
				)
			},

			assertFunc: func(output *ReadOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(validAsaConfig, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validAsaConfig, output)
				}
			},
		},
		{
			testName:  "returns nil ouput when ASA config not found",
			targetUid: asaConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf(urlTemplate, asaConfigUid),
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
			testName:  "return error when fetching ASA Config and remote service encounters issue",
			targetUid: asaConfigUid,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf(urlTemplate, asaConfigUid),
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

			output, err := Read(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), *NewReadInput(asaConfigUid))

			testCase.assertFunc(output, err, t)
		})
	}
}
