package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestReadByUsername(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Should read a user by username", func(t *testing.T) {
		httpmock.Reset()
		expected := user.UserDetails{
			Name:        "barack@example.com",
			ApiOnlyUser: false,
			UserRoles:   []string{"ROLE_ADMIN"},
		}
		httpmock.RegisterResponder(
			"GET",
			"/anubis/rest/v1/users?q=name:"+expected.Name,
			httpmock.NewJsonResponderOrPanic(200, []user.UserDetails{expected}),
		)
		actual, err := user.ReadByUsername(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), user.ReadByUsernameInput{
			Name: expected.Name,
		})
		assert.NotNil(t, actual, "Read output should not be nil")
		assert.Equal(t, *actual, expected)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("Should error if reading a user by username fails", func(t *testing.T) {
		httpmock.Reset()
		name := "donald@example.com"
		httpmock.RegisterResponder(
			"GET",
			"/anubis/rest/v1/users?q=name:"+name,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.ReadByUsername(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), user.ReadByUsernameInput{
			Name: name,
		})
		assert.Nil(t, actual, "Read output should be nil")
		assert.NotNil(t, err, "error should not be nil")
	})
}
