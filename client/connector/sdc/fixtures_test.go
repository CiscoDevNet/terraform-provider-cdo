package sdc_test

import "github.com/CiscoDevNet/terraform-provider-cdo/go-client/connector/sdc"

const (
	cdgUid                   = "11111111-1111-1111-1111-111111111111"
	cdgName                  = "Cloud Connector"
	sdcUid                   = "22222222-2222-2222-2222-222222222222"
	sdcName                  = "My On Prem SDC"
	tenantUid                = "99999999-9999-9999-9999-999999999999"
	tenantName               = "test-tenant-name"
	accessToken              = "test-access-token"
	refreshToken             = "test-refresh-token"
	tokenType                = "test-token-type"
	scope                    = "test-scope"
	baseUrl                  = "https://unittest.cdo.cisco.com"
	host                     = "unittest.cdo.cisco.com"
	serviceConnectivityState = "test-service-connectivity-state"
	state                    = "test-state"
	status                   = "test-status"
)

var (
	bootstrapData = sdc.ComputeBootstrapData(sdcName, accessToken, tenantName, baseUrl, host)
)
