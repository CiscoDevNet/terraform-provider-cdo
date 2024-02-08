package duoadminpanel_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestDuoAdminPanelDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	deleteInput := duoadminpanel.DeleteInput{
		Uid: "test-uid",
	}

	deleteOutput := duoadminpanel.DeleteOutput{}

	baseUrl := "https://test.com"

	testCases := []struct {
		testName   string
		input      duoadminpanel.DeleteInput
		setupFunc  func(input duoadminpanel.DeleteInput)
		assertFunc func(output *duoadminpanel.DeleteOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully delete Duo Admin Panel",
			input:    deleteInput,

			setupFunc: func(input duoadminpanel.DeleteInput) {
				configureDeleteDeviceToReturn(url.DeleteDevice(baseUrl, deleteInput.Uid), deleteOutput)
			},

			assertFunc: func(actualOutput *duoadminpanel.DeleteOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, actualOutput)
				assert.Equal(t, deleteOutput, *actualOutput)
			},
		},
		{
			testName: "fails delete Duo Admin Panel if error",
			input:    deleteInput,

			setupFunc: func(input duoadminpanel.DeleteInput) {
				configureDeleteRequestToError(url.ReadDevice(baseUrl, deleteInput.Uid), "intentional error")
			},

			assertFunc: func(actualOutput *duoadminpanel.DeleteOutput, err error, t *testing.T) {
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

			output, err := duoadminpanel.Delete(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func configureDeleteDeviceToReturn(url string, device duoadminpanel.DeleteOutput) {
	httpmock.RegisterResponder(http.MethodDelete, url, httpmock.NewJsonResponderOrPanic(202, device))
}

func configureDeleteRequestToError(url string, errorBody any) {
	httpmock.RegisterResponder(http.MethodDelete, url, httpmock.NewJsonResponderOrPanic(500, errorBody))
}
