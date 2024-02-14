package tenantsettings

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/settings"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/settings/tenantsettings"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tenantSettingsDataModel struct {
	ID                                    types.String `tfsdk:"id"`
	ChangeRequestSupportEnabled           types.Bool   `tfsdk:"change_request_support_enabled"`
	AutoAcceptDeviceChangesEnabled        types.Bool   `tfsdk:"auto_accept_device_changes_enabled"`
	WebAnalyticsEnabled                   types.Bool   `tfsdk:"web_analytics_enabled"`
	ScheduledDeploymentsEnabled           types.Bool   `tfsdk:"scheduled_deployments_enabled"`
	DenyCiscoSupportAccessToTenantEnabled types.Bool   `tfsdk:"deny_cisco_support_access_to_tenant_enabled"`
	MultiCloudDefenseEnabled              types.Bool   `tfsdk:"multi_cloud_defense_enabled"`
	AutoDiscoverOnPremFmcsEnabled         types.Bool   `tfsdk:"auto_discover_on_prem_fmcs_enabled"`
	ConflictDetectionInterval             types.String `tfsdk:"conflict_detection_interval"`
}

func tenantSettingsDataSourceModelFrom(model settings.TenantSettings) tenantSettingsDataModel {
	return tenantSettingsDataModel{
		ID:                                    types.StringValue(model.Uid.String()),
		ChangeRequestSupportEnabled:           types.BoolValue(model.ChangeRequestSupportEnabled),
		AutoAcceptDeviceChangesEnabled:        types.BoolValue(model.AutoAcceptDeviceChangesEnabled),
		WebAnalyticsEnabled:                   types.BoolValue(model.WebAnalyticsEnabled),
		ScheduledDeploymentsEnabled:           types.BoolValue(model.ScheduledDeploymentsEnabled),
		DenyCiscoSupportAccessToTenantEnabled: types.BoolValue(model.DenyCiscoSupportAccessToTenantEnabled),
		MultiCloudDefenseEnabled:              types.BoolValue(model.MultiCloudDefenseEnabled),
		AutoDiscoverOnPremFmcsEnabled:         types.BoolValue(model.AutoDiscoverOnPremFmcsEnabled),
		ConflictDetectionInterval:             types.StringValue(model.ConflictDetectionInterval.String()),
	}
}

func (model tenantSettingsDataModel) UpdateTenantSettingsInput() tenantsettings.UpdateTenantSettingsInput {
	convertBool := func(tfBool types.Bool) *bool {
		if tfBool.IsNull() || tfBool.IsUnknown() {
			return nil
		}

		return tfBool.ValueBoolPointer()
	}

	convertConflictDetection := func(tfString types.String) *settings.ConflictDetectionInterval {
		if tfString.IsNull() || tfString.IsUnknown() {
			return nil
		}

		interval := settings.ResolveConflictDetectionInterval(tfString.ValueString())
		return &interval
	}

	return tenantsettings.UpdateTenantSettingsInput{
		ChangeRequestSupportEnabled:           convertBool(model.ChangeRequestSupportEnabled),
		AutoAcceptDeviceChangesEnabled:        convertBool(model.AutoAcceptDeviceChangesEnabled),
		WebAnalyticsEnabled:                   convertBool(model.WebAnalyticsEnabled),
		ScheduledDeploymentsEnabled:           convertBool(model.ScheduledDeploymentsEnabled),
		DenyCiscoSupportAccessToTenantEnabled: convertBool(model.DenyCiscoSupportAccessToTenantEnabled),
		MultiCloudDefenseEnabled:              convertBool(model.MultiCloudDefenseEnabled),
		AutoDiscoverOnPremFmcsEnabled:         convertBool(model.AutoDiscoverOnPremFmcsEnabled),
		ConflictDetectionInterval:             convertConflictDetection(model.ConflictDetectionInterval),
	}
}
