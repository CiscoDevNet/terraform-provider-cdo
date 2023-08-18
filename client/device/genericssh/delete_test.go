package genericssh_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/genericssh"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGenericSshDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validGenericSsh := genericssh.DeleteOutput{}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *genericssh.DeleteOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully delete Generic SSH",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodDelete,
					url.DeleteDevice(baseUrl, genericSshUid),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validGenericSsh),
				)
			},
			assertFunc: func(output *genericssh.DeleteOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validGenericSsh, *output)
			},
		},
		{
			testName:  "return error when delete Generic SSH error",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodDelete,
					url.DeleteDevice(baseUrl, genericSshUid),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *genericssh.DeleteOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := genericssh.Delete(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				genericssh.NewDeleteInput(genericSshUid),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
