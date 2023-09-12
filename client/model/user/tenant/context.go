package tenant

type username = string
type userTenantRole = string

type Context struct {
	EulaUsername string               `json:"eulaUsername"`
	Settings     map[username]Setting `json:"settings"`
}

type Setting struct {
	UserTenantRoles []userTenantRole `json:"userTenantRoles"`
	TenantUid       string           `json:"tenantUid"`
	Username        string           `json:"username"`
}
