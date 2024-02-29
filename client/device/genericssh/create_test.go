package genericssh_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/genericssh"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	internalTesting "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGenericSshCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validGenericSsh := genericssh.CreateOutput{
		Uid:   genericSshUid,
		Name:  genericSshName,
		State: state.DONE,
		Tags:  internalTesting.NewTestingTags(),
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *genericssh.CreateOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully create Generic SSH",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					url.CreateDevice(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validGenericSsh),
				)
			},
			assertFunc: func(output *genericssh.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validGenericSsh, *output)
			},
		},
		{
			testName:  "return error when creating Generic SSH error",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPost,
					url.CreateDevice(baseUrl),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *genericssh.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := genericssh.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				genericssh.NewCreateInput(genericSshUid, genericSshConnectorUid, genericSshConnectorSocketAddress, validGenericSsh.Tags),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
