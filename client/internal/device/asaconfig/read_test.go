package asaconfig

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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
				assert.Nil(t, err)
				assert.NotNil(t, output)

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
				assert.Nil(t, output)
				assert.NotNil(t, err)
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
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := Read(context.Background(), *http.MustNewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), *NewReadInput(asaConfigUid))

			testCase.assertFunc(output, err, t)
		})
	}
}
