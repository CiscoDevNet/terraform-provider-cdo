package duoadminpanel_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/jarcoal/httpmock"
)

func TestDuoAdminPanelReadByUid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	readInput := duoadminpanel.ReadByUidInput{
		Uid: "test-uid",
	}

	readDevice := duoadminpanel.ReadOutput{
		Uid:  "test-uid",
		Name: "test-name",
	}

	baseUrl := "https://test.com"

	testCases := []struct {
		testName   string
		input      duoadminpanel.ReadByUidInput
		setupFunc  func(input duoadminpanel.ReadByUidInput)
		assertFunc func(output *duoadminpanel.ReadOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully read Duo Admin Panel",
			input:    readInput,

			setupFunc: func(input duoadminpanel.ReadByUidInput) {
				internalTesting.MockGetOk(url.ReadDevice(baseUrl, readDevice.Uid), readDevice)
			},

			assertFunc: func(actualOutput *duoadminpanel.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, actualOutput)
				assert.Equal(t, readDevice, *actualOutput)
			},
		},
		{
			testName: "fails read Duo Admin Panel if error",
			input:    readInput,

			setupFunc: func(input duoadminpanel.ReadByUidInput) {
				internalTesting.MockGetError(url.ReadDevice(baseUrl, readDevice.Uid), "intentional error")
			},

			assertFunc: func(actualOutput *duoadminpanel.CreateOutput, err error, t *testing.T) {
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

			output, err := duoadminpanel.ReadByUid(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
