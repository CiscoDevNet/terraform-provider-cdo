package connectoronboarding_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/connectoronboarding"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/status"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	activeConnector := connector.ReadOutput{
		Uid:              "",
		Name:             "test-sdc",
		DefaultConnector: false,
		Cdg:              false,
		TenantUid:        "",
		PublicKey:        model.PublicKey{},
		ConnectorStatus:  status.Active,
	}

	onboardingConnector := connector.ReadOutput{
		Uid:              "",
		Name:             "test-sdc",
		DefaultConnector: false,
		Cdg:              false,
		TenantUid:        "",
		PublicKey:        model.PublicKey{},
		ConnectorStatus:  status.Onboarding,
	}

	baseUrl := "https://unittest.cdo.cisco.com"

	testCases := []struct {
		testName   string
		input      connectoronboarding.CreateInput
		setupFunc  func()
		assertFunc func(output *connectoronboarding.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "should finish on connector Active status",
			input:    connectoronboarding.NewCreateInput(activeConnector.Name),
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadConnectorByName(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, []connector.ReadOutput{activeConnector}),
				)
			},
			assertFunc: func(output *connectoronboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, output)
				assert.Nil(t, err)
				assert.Equal(t, *output, activeConnector)
			},
		},
		{
			testName: "should ends on connector Active status after an Onboarding state",
			input:    connectoronboarding.NewCreateInput(activeConnector.Name),
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadConnectorByName(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, []connector.ReadOutput{onboardingConnector}),
				)
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadConnectorByName(baseUrl),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, []connector.ReadOutput{activeConnector}),
				)
			},
			assertFunc: func(output *connectoronboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, output)
				assert.Nil(t, err)
				assert.Equal(t, *output, activeConnector)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := connectoronboarding.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
