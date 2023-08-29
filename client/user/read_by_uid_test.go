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

func TestReadByUid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Should read a user by UID", func(t *testing.T) {
		httpmock.Reset()
		uid := "sample-uid"
		expected := user.UserDetails{
			Name:        "dubya@example.com",
			ApiOnlyUser: false,
			UserRoles:   []string{"ROLE_SUPER_ADMIN"},
		}
		httpmock.RegisterResponder(
			"GET",
			"/anubis/rest/v1/users/"+uid,
			httpmock.NewJsonResponderOrPanic(200, expected),
		)
		actual, err := user.ReadByUid(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), user.ReadByUidInput{
			Uid: uid,
		})
		assert.NotNil(t, actual, "Read output should not be nil")
		assert.Equal(t, *actual, expected)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("Should error if reading a user by UID fails", func(t *testing.T) {
		httpmock.Reset()
		uid := "sample-user-uid"
		httpmock.RegisterResponder(
			"GET",
			"/anubis/rest/v1/users/"+uid,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.ReadByUid(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), user.ReadByUidInput{
			Uid: uid,
		})
		assert.Nil(t, actual, "Read output should be nil")
		assert.NotNil(t, err, "error should not be nil")
	})
}
