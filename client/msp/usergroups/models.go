package usergroups

type MspManagedUserGroupInput struct {
	GroupIdentifier string  `json:"groupIdentifier"`
	IssuerUrl       string  `json:"issuerUrl"`
	Name            string  `json:"name"`
	Role            string  `json:"role"`
	Notes           *string `json:"notes"`
}

type MspManagedUserGroup struct {
	GroupIdentifier string  `json:"groupIdentifier"`
	IssuerUrl       string  `json:"issuerUrl"`
	Name            string  `json:"name"`
	Role            string  `json:"role"`
	Notes           *string `json:"notes"`
	Uid             string  `json:"uid"`
}

type MspManagedUserGroupPage struct {
	Count  int                   `json:"count"`
	Offset int                   `json:"offset"`
	Limit  int                   `json:"limit"`
	Items  []MspManagedUserGroup `json:"items"`
}

type MspManagedUserGroupDeleteInput struct {
	UserGroupUids []string `json:"userGroupUids"`
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}
