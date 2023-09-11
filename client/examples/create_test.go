package examples_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/examples"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
	"time"
)

func TestExampleCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName   string
		input      examples.CreateInput
		setupFunc  func()
		assertFunc func(output *examples.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "example test",
			input:    examples.NewCreateInput("unittest-device-uid"),
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadDevice("https://unittest.cdo.cisco.com", "unittest-device-uid"),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, "{\"a\":\"b\"}"),
				)
			},
			assertFunc: func(output *examples.CreateOutput, err error, t *testing.T) {
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := examples.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig("https://unittest.cdo.cisco.com", "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
