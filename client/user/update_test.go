package user_test

import (
	"context"
	netHttp "net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
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
	t.Run("Successfully update an existing user", func(t *testing.T) {
		httpmock.Reset()
		expected := model.UserDetails{
			Name:        "jfk@example.com",
			ApiOnlyUser: false,
			UserRoles:   []string{"ROLE_SUPER_ADMIN"},
		}

		httpmock.RegisterResponder(
			netHttp.MethodPut,
			"/anubis/rest/v1/users/"+validUserTenantAssociation.Source.Uid,
			httpmock.NewJsonResponderOrPanic(200, validUserTenantAssociation),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/users/"+validUserTenantAssociation.Source.Uid,
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		actual, err := user.Update(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewUpdateByUidInput(validUserTenantAssociation.Source.Uid, []string{"ROLE_SUPER_ADMIN"}))

		assert.NotNil(t, actual, "User details returned must not be nil")
		assert.Equal(t, expected, *actual, "Actual user details do not match expected")
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("should error if failed to update user", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			netHttp.MethodPut,
			"/anubis/rest/v1/users/"+validUserTenantAssociation.Source.Uid,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.Update(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewUpdateByUidInput(validUserTenantAssociation.Source.Uid, []string{"ROLE_SUPER_ADMIN"}))

		assert.Nil(t, actual, "Expected actual user not to be updated")
		assert.NotNil(t, err, "Expected error")
	})

	t.Run("should error if failed to read user details after update", func(t *testing.T) {
		httpmock.RegisterResponder(
			netHttp.MethodPut,
			"/anubis/rest/v1/users/"+validUserTenantAssociation.Source.Uid,
			httpmock.NewJsonResponderOrPanic(200, validUserTenantAssociation),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/users/"+validUserTenantAssociation.Source.Uid,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.Update(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewUpdateByUidInput(validUserTenantAssociation.Source.Uid, []string{"ROLE_SUPER_ADMIN"}))
		assert.Nil(t, actual, "Expected actual user not to be updated")
		assert.NotNil(t, err, "Expected error")
	})
}
