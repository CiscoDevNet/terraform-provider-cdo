package sec_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sec"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestSecCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	baseUrl := "https://unittest.cdo.cisco.com"
	domain := "unittest.cdo.cisco.com"
	successTokenResponse := user.NewGetTokenOutputBuilder().
		AccessToken("test-access-token").
		TenantName("test-tenant-name").
		Build()
	successCreateOutput := sec.NewCreateOutputBuilder().
		Uid("test-sec-uid").
		SecBootstrapData("sec-test-bootstrap-data").
		Name("test-sec-name").
		CdoBoostrapData(sec.ComputeEventOnlyBootstrapData(successTokenResponse.AccessToken, successTokenResponse.TenantName, baseUrl, domain)).
		Build()
	successReadOutput := sec.NewReadOutputBuilder().
		Uid(successCreateOutput.Uid).
		Name(successCreateOutput.Name).
		BootStrapData(successCreateOutput.SecBootstrapData).
		TokenExpiryTime(123).
		Build()
	noBootstrapDataReadOutput := sec.NewReadOutputBuilder().
		Uid(successCreateOutput.Uid).
		Name(successCreateOutput.Name).
		BootStrapData("").
		TokenExpiryTime(123).
		Build()

	testCases := []struct {
		testName   string
		input      sec.CreateInput
		setupFunc  func()
		assertFunc func(output *sec.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "Create successfully",
			input:    sec.NewCreateInputBuilder().Build(),
			setupFunc: func() {
				startedCreateSecDevice(baseUrl, successCreateOutput)
				createSecStateMachineEnded(baseUrl, "eventingPushRequest", state.DONE)
				secBootstrapDataGenerated(baseUrl, successReadOutput)
				cdoBootstrapDataGenerated(baseUrl, successTokenResponse)
				secNameIsUpdated(baseUrl, successReadOutput)
			},
			assertFunc: func(output *sec.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.Equal(t, successCreateOutput, *output)
			},
		},
		{
			testName: "Error when failed to initiate creation of SEC device in CDO",
			input:    sec.NewCreateInputBuilder().Build(),
			setupFunc: func() {
				startedCreateSecDevice(baseUrl, successCreateOutput)
				SecStateMachineEndedInFailure(baseUrl)
				secBootstrapDataGenerated(baseUrl, successReadOutput)
				cdoBootstrapDataGenerated(baseUrl, successTokenResponse)
				secNameIsUpdated(baseUrl, successReadOutput)
			},
			assertFunc: func(output *sec.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, url.ReadStateMachineInstance(baseUrl))
			},
		},
		{
			testName: "Error when failed to create SEC device in CDO",
			input:    sec.NewCreateInputBuilder().Build(),
			setupFunc: func() {
				failedToStartCreateSecDevice(baseUrl)
				createSecStateMachineEnded(baseUrl, "eventingPushRequest", state.DONE)
				secBootstrapDataGenerated(baseUrl, successReadOutput)
				cdoBootstrapDataGenerated(baseUrl, successTokenResponse)
				secNameIsUpdated(baseUrl, successReadOutput)
			},
			assertFunc: func(output *sec.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, url.CreateSec(baseUrl))
			},
		},
		{
			testName: "Error when SEC bootstrap data failed to generate",
			input:    sec.NewCreateInputBuilder().Build(),
			setupFunc: func() {
				startedCreateSecDevice(baseUrl, successCreateOutput)
				createSecStateMachineEnded(baseUrl, "eventingPushRequest", state.DONE)
				secBootstrapDataGenerated(baseUrl, noBootstrapDataReadOutput)
				cdoBootstrapDataGenerated(baseUrl, successTokenResponse)
				// getting the bootstrap data and checking if the SEC name is updated are the same call, so skipping secNameIsUpdated(baseUrl, successReadOutput) so that it does not overwrite
			},
			assertFunc: func(output *sec.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, "SEC bootstrap data not found")
			},
		},
		{
			testName: "Error when fail to get SEC bootstrap data",
			input:    sec.NewCreateInputBuilder().Build(),
			setupFunc: func() {
				startedCreateSecDevice(baseUrl, successCreateOutput)
				createSecStateMachineEnded(baseUrl, "eventingPushRequest", state.DONE)
				failedToGeneratedSecBootstrapData(baseUrl, successReadOutput.Uid)
				cdoBootstrapDataGenerated(baseUrl, successTokenResponse)
				// getting the bootstrap data and checking if the SEC name is updated are the same call, so skipping secNameIsUpdated(baseUrl, successReadOutput) so that it does not overwrite
			},
			assertFunc: func(output *sec.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, url.ReadSec(baseUrl, successReadOutput.Uid))
			},
		},
		{
			testName: "Error when failed to generate CDO bootstrap data",
			input:    sec.NewCreateInputBuilder().Build(),
			setupFunc: func() {
				startedCreateSecDevice(baseUrl, successCreateOutput)
				createSecStateMachineEnded(baseUrl, "eventingPushRequest", state.DONE)
				secBootstrapDataGenerated(baseUrl, successReadOutput)
				failedToGenerateCdoBootstrapData(baseUrl)
				secNameIsUpdated(baseUrl, successReadOutput)
			},
			assertFunc: func(output *sec.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, url.UserToken(baseUrl))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := sec.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func startedCreateSecDevice(baseUrl string, response sec.CreateOutput) {
	httpmock.RegisterResponder(
		http.MethodPost,
		url.CreateSec(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, response),
	)
}

func createSecStateMachineEnded(baseUrl string, stateMachineIdentifier string, state state.Type) {
	response := statemachine.NewReadInstanceByDeviceUidOutputBuilder().
		StateMachineInstanceCondition(state).
		StateMachineIdentifier(stateMachineIdentifier).
		Build()
	// check started successfully
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadStateMachineInstance(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, []statemachine.ReadInstanceByDeviceUidOutput{response}),
	)
	// check ended successfully
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadStateMachineInstance(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, []statemachine.ReadInstanceByDeviceUidOutput{response}),
	)
}

func secBootstrapDataGenerated(baseUrl string, response sec.ReadOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadSec(baseUrl, response.Uid),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, response),
	)
}

func cdoBootstrapDataGenerated(baseUrl string, response user.GetTokenOutput) {
	httpmock.RegisterResponder(
		http.MethodPost,
		url.UserToken(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, response),
	)
}

func secNameIsUpdated(baseUrl string, response sec.ReadOutput) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadSec(baseUrl, response.Uid),
		httpmock.NewJsonResponderOrPanic(http.StatusOK, response),
	)
}

func failedToStartCreateSecDevice(baseUrl string) {
	httpmock.RegisterResponder(
		http.MethodPost,
		url.CreateSec(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "intentional error"),
	)
}

func SecStateMachineEndedInFailure(baseUrl string) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadStateMachineInstance(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "intentional error"),
	)
}

func failedToGeneratedSecBootstrapData(baseUrl string, uid string) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadSec(baseUrl, uid),
		httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "intentional error"),
	)
}

func failedToGenerateCdoBootstrapData(baseUrl string) {
	httpmock.RegisterResponder(
		http.MethodPost,
		url.UserToken(baseUrl),
		httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "intentional error"),
	)
}

func secNameFailedToUpdated(baseUrl string, uid string) {
	httpmock.RegisterResponder(
		http.MethodGet,
		url.ReadSec(baseUrl, uid),
		httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, "intentional error"),
	)
}
