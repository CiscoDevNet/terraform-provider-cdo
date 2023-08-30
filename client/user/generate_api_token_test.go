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

func TestGenerateApiToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "lbj@example.com"
	t.Run("Successfully generate an API token", func(t *testing.T) {
		apiTokenResponse := user.ApiTokenResponse{ApiToken: "jwt-token"}
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/anubis/rest/v1/oauth/token/"+username,
			httpmock.NewJsonResponderOrPanic(200, apiTokenResponse),
		)

		actual, err := user.GenerateApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewGenerateApiTokenInput(username))
		assert.NotNil(t, actual, "API token response should not be nil")
		assert.Equal(t, apiTokenResponse, *actual)
		assert.Nil(t, err, "Error cannot be non-nil")
	})

	t.Run("Should fail if API token generation failed", func(t *testing.T) {
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/anubis/rest/v1/oauth/token/"+username,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.GenerateApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewGenerateApiTokenInput(username))
		assert.Nil(t, actual, "API token response should be nil")
		assert.NotNil(t, err, "Error cannot be nil")
	})
}
