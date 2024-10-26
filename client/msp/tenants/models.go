package tenants

type MspCreateTenantInput struct {
	Name        string `json:"tenantName"`
	DisplayName string `json:"displayName"`
}

type MspTenantOutput struct {
	Uid         string `json:"uid"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Region      string `json:"region"`
}

type MspTenantsOutput struct {
	Count  int               `json:"count"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
	Items  []MspTenantOutput `json:"items"`
}

type ReadByUidInput struct {
	Uid string `json:"uid"`
}

type ReadByNameInput struct {
	Name string `json:"name"`
}
