package users

type MspUsersInput struct {
	TenantUid string        `json:"tenantUid"`
	Users     []UserDetails `json:"users"`
}

type MspUsersPublicApiInput struct {
	TenantUid string                      `json:"tenantUid"`
	Users     []UserDetailsPublicApiInput `json:"users"`
}

type UserDetailsPublicApiInput struct {
	Uid         string `json:"uid"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	ApiOnlyUser bool   `json:"apiOnlyUser"`
}

type MspDeleteUsersInput struct {
	TenantUid string   `json:"tenantUid"`
	Usernames []string `json:"usernames"`
}

type UserDetails struct {
	Uid         string   `json:"uid"`
	Username    string   `json:"name"`
	Roles       []string `json:"roles"`
	ApiOnlyUser bool     `json:"apiOnlyUser"`
}

type UserPage struct {
	Count  int           `json:"count"`
	Offset int           `json:"offset"`
	Limit  int           `json:"limit"`
	Items  []UserDetails `json:"items"`
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}
