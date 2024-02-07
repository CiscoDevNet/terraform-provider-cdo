package duoadminpanel_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestDuoAdminPanelUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	updateInput := duoadminpanel.UpdateInput{
		Uid:  "test-uid",
		Name: "test-name",
		Tags: tags.New([]string{"1", "2", "3"}...),
	}

	updateOutput := duoadminpanel.UpdateOutput{
		Uid:  updateInput.Uid,
		Name: updateInput.Name,
	}

	baseUrl := "https://test.com"

	testCases := []struct {
		testName   string
		input      duoadminpanel.UpdateInput
		setupFunc  func(input duoadminpanel.UpdateInput)
		assertFunc func(output *duoadminpanel.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully update Duo Admin Panel",
			input:    updateInput,

			setupFunc: func(input duoadminpanel.UpdateInput) {
				configureUpdateToReturn(url.ReadDevice(baseUrl, updateOutput.Uid), updateOutput)
			},

			assertFunc: func(actualOutput *duoadminpanel.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, actualOutput)
				assert.Equal(t, updateOutput, *actualOutput)
			},
		},
		{
			testName: "fails update Duo Admin Panel if error",
			input:    updateInput,

			setupFunc: func(input duoadminpanel.UpdateInput) {
				configurePutRequestToError(url.UpdateDevice(baseUrl, updateOutput.Uid), "intentional error")
			},

			assertFunc: func(actualOutput *duoadminpanel.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, actualOutput)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, "intentional error")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := duoadminpanel.Update(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func configureUpdateToReturn(url string, device duoadminpanel.UpdateOutput) {
	httpmock.RegisterResponder(http.MethodPut, url, httpmock.NewJsonResponderOrPanic(200, device))
}

func configurePutRequestToError(url string, errorBody any) {
	httpmock.RegisterResponder(http.MethodPut, url, httpmock.NewJsonResponderOrPanic(500, errorBody))
}
