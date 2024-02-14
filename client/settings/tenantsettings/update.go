package tenantsettings

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/settings"
)

type UpdateTenantSettingsInput struct {
	ChangeRequestSupportEnabled           *bool                               `json:"changeRequestSupport,omitempty"`
	AutoAcceptDeviceChangesEnabled        *bool                               `json:"autoAcceptDeviceChanges,omitempty"`
	WebAnalyticsEnabled                   *bool                               `json:"webAnalytics,omitempty"`
	ScheduledDeploymentsEnabled           *bool                               `json:"scheduledDeployments,omitempty"`
	DenyCiscoSupportAccessToTenantEnabled *bool                               `json:"denyCiscoSupportAccessToTenant,omitempty"`
	MultiCloudDefenseEnabled              *bool                               `json:"multicloudDefense,omitempty"`
	AutoDiscoverOnPremFmcsEnabled         *bool                               `json:"autoDiscoverOnPremFmcs,omitempty"`
	ConflictDetectionInterval             *settings.ConflictDetectionInterval `json:"conflictDetectionInterval,omitempty"`
}

func Update(ctx context.Context, client http.Client, input UpdateTenantSettingsInput) (*settings.TenantSettings, error) {
	updateURL := url.UpdateTenantSettings(client.BaseUrl())
	req := client.NewPatch(ctx, updateURL, input)

	var settings settings.TenantSettings
	if err := req.Send(&settings); err != nil {
		return nil, err
	}

	return &settings, nil
}
