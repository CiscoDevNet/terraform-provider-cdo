package users

type MspCreateUsersInput struct {
	TenantUid string        `json:"tenantUid"`
	Users     []UserDetails `json:"users"`
}

type MspDeleteUsersInput struct {
	TenantUid string   `json:"tenantUid"`
	Usernames []string `json:"usernames"`
}

type UserDetails struct {
	Username    string `json:"username"`
	Role        string `json:"role"`
	ApiOnlyUser bool   `json:"apiOnlyUser"`
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}
