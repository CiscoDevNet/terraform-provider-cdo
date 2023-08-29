package user_test

import (
	"context"
	"fmt"
	netHttp "net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validUserTenantAssociation := user.UserTenantAssociation{
		Uid: "association-uuid",
		Source: user.Association{
			Uid:       "sample-uuid",
			Namespace: "systemdb",
			Type:      "users",
		},
	}

	t.Run("successfully create user", func(t *testing.T) {
		httpmock.Reset()
		expected := user.UserDetails{
			Name:        "george@example.com",
			ApiOnlyUser: false,
			UserRoles:   []string{"ROLE_SUPER_ADMIN"},
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/anubis/rest/v1/users/%s", expected.Name),
			httpmock.NewJsonResponderOrPanic(200, validUserTenantAssociation),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/users/"+validUserTenantAssociation.Source.Uid,
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		actual, err := user.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewCreateUserInput(expected.Name, expected.UserRoles[0], expected.ApiOnlyUser))
		assert.NotNil(t, actual, "User details returned must not be nil")
		assert.Equal(t, expected, *actual, "Actual user details do not match expected")
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("should error if failed to create user", func(t *testing.T) {
		username := "bill@example.com"
		apiOnlyUser := false
		userRoles := []string{"ROLE_READ_ONLY"}
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/anubis/rest/v1/users/%s", username),
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewCreateUserInput(username, userRoles[0], apiOnlyUser))
		assert.Nil(t, actual, "Expected actual user not to be created")
		assert.NotNil(t, err, "Expected error")
	})

	t.Run("should error if failed to read user details after creation", func(t *testing.T) {
		username := "bill@example.com"
		apiOnlyUser := false
		userRoles := []string{"ROLE_READ_ONLY"}
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/anubis/rest/v1/users/%s", username),
			httpmock.NewJsonResponderOrPanic(200, validUserTenantAssociation),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/users/"+validUserTenantAssociation.Source.Uid,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewCreateUserInput(username, userRoles[0], apiOnlyUser))
		assert.Nil(t, actual, "Expected actual user not to be created")
		assert.NotNil(t, err, "Expected error")
	})
}
