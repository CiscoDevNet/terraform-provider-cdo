package genericssh_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/genericssh"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGenericSshRead(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validGenericSsh := genericssh.ReadOutput{
		Uid:   genericSshUid,
		Name:  genericSshName,
		State: state.DONE,
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *genericssh.ReadOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully read Generic SSH",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadDevice(baseUrl, genericSshUid),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validGenericSsh),
				)
			},
			assertFunc: func(output *genericssh.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validGenericSsh, *output)
			},
		},
		{
			testName:  "returns nil output when Generic SSH not found",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadDevice(baseUrl, genericSshUid),
					httpmock.NewStringResponder(http.StatusNotFound, ""),
				)
			},
			assertFunc: func(output *genericssh.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
		{
			testName:  "return error when reading Generic SSH error",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadDevice(baseUrl, genericSshUid),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *genericssh.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := genericssh.Read(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				*genericssh.NewReadInput(genericSshUid),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
