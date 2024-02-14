package tenantsettings_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/settings"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/settings/tenantsettings"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTenantSettings(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	baseURL := "https://unit-test.cdo.cisco.com"
	tenantSettings := settings.TenantSettings{
		Uid:                                   uuid.New(),
		ChangeRequestSupportEnabled:           false,
		AutoAcceptDeviceChangesEnabled:        false,
		WebAnalyticsEnabled:                   false,
		ScheduledDeploymentsEnabled:           false,
		DenyCiscoSupportAccessToTenantEnabled: false,
		MultiCloudDefenseEnabled:              false,
		AutoDiscoverOnPremFmcsEnabled:         false,
		ConflictDetectionInterval:             settings.ConflictDetectionIntervalEvery10Minutes,
	}
	tenMinsConflictDetection := settings.ConflictDetectionIntervalEvery10Minutes
	updateInput := tenantsettings.UpdateTenantSettingsInput{
		ConflictDetectionInterval: &tenMinsConflictDetection,
	}

	testCases := []struct {
		testName    string
		updateInput tenantsettings.UpdateTenantSettingsInput
		setupFunc   func()
		assertFunc  func(output *settings.TenantSettings, err error, t *testing.T)
	}{
		{
			testName:    "successfully update tenant settings",
			updateInput: updateInput,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPatch,
					url.ReadTenantSettings(baseURL),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, tenantSettings),
				)
			},
			assertFunc: func(output *settings.TenantSettings, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, tenantSettings, *output)
			},
		},
		{
			testName:    "return error when update tenant settings error",
			updateInput: updateInput,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadTenantSettings(baseURL),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *settings.TenantSettings, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := tenantsettings.Update(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseURL, "a_valid_token", 0, 0, time.Minute),
				testCase.updateInput,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
