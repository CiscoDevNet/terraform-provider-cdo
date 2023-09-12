package tenant

type TenantDetailsDetails struct {
	TenantName             string `json:"TenantName"`
	TenantOrganizationName string `json:"TenantOrganizationName"`
	TenantPayType          string `json:"TenantPayType"`
	TenantUid              string `json:"TenantUid"`
}

type UserAuthentication struct {
	Details TenantDetailsDetails `json:"details"`
}

type ReadTenantDetailsOutput struct {
	UserAuthentication UserAuthentication `json:"userAuthentication"`
}
