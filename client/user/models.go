package user

type CreateUserInput struct {
	Username    string
	UserRoles   string
	ApiOnlyUser bool
}

type UpdateUserInput struct {
	Uid       string
	UserRoles []string
}

type DeleteUserInput struct {
	Uid string
}

type DeleteUserOutput struct{}

// CreateUser endpoint returns a user-tenant association for whatever reason
type UserTenantAssociation struct {
	Uid    string      `json:"uid"`
	Source Association `json:"source"`
}

type UserDetails struct {
	Uid                 string   `json:"uid"`
	Name                string   `json:"name"`
	UserRoles           []string `json:"roles"`
	ApiOnlyUser         bool     `json:"isApiOnlyUser"`
	LastSuccessfulLogin int64    `json:"lastSuccessfulLogin"`
}

type Association struct {
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Uid       string `json:"uid"`
}

type ReadByUsernameInput struct {
	Name string `json:"name"`
}

type ReadByUidInput struct {
	Uid string `json:"uid"`
}

func NewCreateUserInput(username string, userRoles string, apiOnlyUser bool) *CreateUserInput {
	return &CreateUserInput{
		Username:    username,
		UserRoles:   userRoles,
		ApiOnlyUser: apiOnlyUser,
	}
}

func NewReadByUsernameInput(name string) *ReadByUsernameInput {
	return &ReadByUsernameInput{
		Name: name,
	}
}

func NewReadByUidInput(uid string) *ReadByUidInput {
	return &ReadByUidInput{
		Uid: uid,
	}
}

func NewUpdateByUidInput(uid string, userRoles []string) *UpdateUserInput {
	return &UpdateUserInput{
		Uid:       uid,
		UserRoles: userRoles,
	}
}
