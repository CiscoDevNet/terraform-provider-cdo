package settings

import "github.com/google/uuid"

type TenantSettings struct {
	Uid                                   uuid.UUID                 `json:"uid"`
	ChangeRequestSupportEnabled           bool                      `json:"changeRequestSupport"`
	AutoAcceptDeviceChangesEnabled        bool                      `json:"autoAcceptDeviceChanges"`
	WebAnalyticsEnabled                   bool                      `json:"webAnalytics"`
	ScheduledDeploymentsEnabled           bool                      `json:"scheduledDeployments"`
	DenyCiscoSupportAccessToTenantEnabled bool                      `json:"denyCiscoSupportAccessToTenant"`
	MultiCloudDefenseEnabled              bool                      `json:"multicloudDefense"`
	AutoDiscoverOnPremFmcsEnabled         bool                      `json:"autoDiscoverOnPremFmcs"`
	ConflictDetectionInterval             ConflictDetectionInterval `json:"conflictDetectionInterval"`
}
