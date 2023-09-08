package auth

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth/role"

type Info struct {
	UserAuthentication Authentication `json:"userAuthentication"`
}

type Authority struct {
	Authority role.Type `json:"authority"`
}

type Authentication struct {
	Authorities   []Authority `json:"authorities"`
	Details       Details     `json:"details"`
	Authenticated bool        `json:"authenticated"`
	Principle     string      `json:"principal"`
	Name          string      `json:"name"`
}

type Details struct {
	TenantUid              string `json:"TenantUid"`
	TenantName             string `json:"TenantName"`
	SseTenantUid           string `json:"sseTenantUid"`
	TenantOrganizationName string `json:"TenantOrganizationName"`
	TenantDbFeatures       string `json:"TenantDbFeatures"`
	TenantUserRoles        string `json:"TenantUserRoles"`
	TenantDatabaseName     string `json:"TenantDatabaseName"`
	TenantPayType          string `json:"TenantPayType"`
}
