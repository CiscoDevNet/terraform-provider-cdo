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

type ReadByUidInput struct {
	Uid string `json:"uid"`
}
