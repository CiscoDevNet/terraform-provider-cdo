package user

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
)

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

type RevokeApiTokenOutput struct{}

// CreateUser endpoint returns a user-tenant association for whatever reason
type UserTenantAssociation struct {
	Uid    string      `json:"uid"`
	Source Association `json:"source"`
}

type CreateUserOutput = model.UserDetails

type UpdateUserOutput = model.UserDetails

type ReadUserOutput = model.UserDetails

type Association struct {
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Uid       string `json:"uid"`
}

type ReadByUsernameInput struct {
	Name string `json:"name"`
}

type GenerateApiTokenInput struct {
	Name string `json:"name"`
}

type RevokeApiTokenInput struct {
	Name string `json:"name"`
}

type RevokeOAuthTokenInput struct {
	ApiTokenId string `json:"apiTokenId"`
}

type ReadByUidInput struct {
	Uid string `json:"uid"`
}

type ApiTokenResponse struct {
	ApiToken string `json:"access_token"`
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

func NewGenerateApiTokenInput(name string) *GenerateApiTokenInput {
	return &GenerateApiTokenInput{
		Name: name,
	}
}

func NewRevokeApiTokenInput(name string) *RevokeApiTokenInput {
	return &RevokeApiTokenInput{
		Name: name,
	}
}

func NewRevokeOAuthTokenInput(apiTokenId string) *RevokeOAuthTokenInput {
	return &RevokeOAuthTokenInput{
		ApiTokenId: apiTokenId,
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
