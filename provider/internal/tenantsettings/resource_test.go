package tenantsettings_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/settings"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testTenantSettingsResourceConfigTemplate = `resource "cdo_tenant_settings" "test" {
	change_request_support_enabled              = {{.ChangeRequestSupportEnabled}}
	auto_accept_device_changes_enabled          = {{.AutoAcceptDeviceChangesEnabled}}
	web_analytics_enabled                       = {{.WebAnalyticsEnabled}}
	scheduled_deployments_enabled               = {{.ScheduledDeploymentsEnabled}}
	deny_cisco_support_access_to_tenant_enabled = {{.DenyCiscoSupportAccessToTenantEnabled}}
	multi_cloud_defense_enabled                 = {{.MultiCloudDefenseEnabled}}
	auto_discover_on_prem_fmcs_enabled          = {{.AutoDiscoverOnPremFmcsEnabled}}
	conflict_detection_interval                 = "{{.ConflictDetectionInterval}}"
}`

var createTenantSettingsResourceConfig = acctest.MustParseTemplate(testTenantSettingsResourceConfigTemplate, settings.TenantSettings{
	ChangeRequestSupportEnabled:           true,
	AutoAcceptDeviceChangesEnabled:        true,
	WebAnalyticsEnabled:                   true,
	ScheduledDeploymentsEnabled:           true,
	DenyCiscoSupportAccessToTenantEnabled: true,
	MultiCloudDefenseEnabled:              true,
	AutoDiscoverOnPremFmcsEnabled:         true,
	ConflictDetectionInterval:             settings.ConflictDetectionIntervalEvery6Hours,
})

var updateTenantSettingsResourceConfig = acctest.MustParseTemplate(testTenantSettingsResourceConfigTemplate, settings.TenantSettings{
	ChangeRequestSupportEnabled:           false,
	AutoAcceptDeviceChangesEnabled:        false,
	WebAnalyticsEnabled:                   false,
	ScheduledDeploymentsEnabled:           false,
	DenyCiscoSupportAccessToTenantEnabled: false,
	MultiCloudDefenseEnabled:              false,
	AutoDiscoverOnPremFmcsEnabled:         false,
	ConflictDetectionInterval:             settings.ConflictDetectionIntervalEvery24Hours,
})

func TestAccTenantSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + createTenantSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "id", acctest.Env.TenantSettingsTenantUid()),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "change_request_support_enabled", "true"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "auto_accept_device_changes_enabled", "true"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "web_analytics_enabled", "true"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "scheduled_deployments_enabled", "true"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "deny_cisco_support_access_to_tenant_enabled", "true"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "multi_cloud_defense_enabled", "true"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "auto_discover_on_prem_fmcs_enabled", "true"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "conflict_detection_interval", string(settings.ConflictDetectionIntervalEvery6Hours)),
				),
			},

			// Update testing
			{
				Config: acctest.ProviderConfig() + updateTenantSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "id", acctest.Env.TenantSettingsTenantUid()),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "change_request_support_enabled", "false"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "auto_accept_device_changes_enabled", "false"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "web_analytics_enabled", "false"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "scheduled_deployments_enabled", "false"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "deny_cisco_support_access_to_tenant_enabled", "false"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "multi_cloud_defense_enabled", "false"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "auto_discover_on_prem_fmcs_enabled", "false"),
					resource.TestCheckResourceAttr("cdo_tenant_settings.test", "conflict_detection_interval", string(settings.ConflictDetectionIntervalEvery24Hours)),
				),
			},
		},
	})
}
