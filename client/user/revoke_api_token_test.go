package user_test

import (
	"context"
	netHttp "net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestRevokeApiToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	userDetails := user.UserDetails{
		Name:        "barack@example.com",
		ApiOnlyUser: false,
		UserRoles:   []string{"ROLE_ADMIN"},
		ApiTokenId:  "api-token-id",
	}

	t.Run("Successfully revoke an API token", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/users?q=name:"+userDetails.Name,
			httpmock.NewJsonResponderOrPanic(200, []user.UserDetails{userDetails}),
		)

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/anubis/rest/v1/oauth/revoke/"+userDetails.ApiTokenId,
			httpmock.NewJsonResponderOrPanic(200, nil),
		)

		actual, err := user.RevokeApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewRevokeApiTokenInput(userDetails.Name))
		assert.NotNil(t, actual, "Revocation response should not be nil")
		assert.Nil(t, err, "Error cannot be non-nil")
	})

	t.Run("Should fail if API token revocation failed", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/users?q=name:"+userDetails.Name,
			httpmock.NewJsonResponderOrPanic(200, []user.UserDetails{userDetails}),
		)

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/anubis/rest/v1/oauth/revoke/"+userDetails.ApiTokenId,
			httpmock.NewJsonResponderOrPanic(400, nil),
		)

		actual, err := user.RevokeApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewRevokeApiTokenInput(userDetails.Name))
		assert.Nil(t, actual, "Revocation response should be nil")
		assert.NotNil(t, err, "Error should be nil")
	})

	t.Run("Should fail if API token revocation failed because the user could not be found", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/users?q=name:"+userDetails.Name,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.RevokeApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewRevokeApiTokenInput(userDetails.Name))
		assert.Nil(t, actual, "Revocation response should be nil")
		assert.NotNil(t, err, "Error should be nil")
	})
}
