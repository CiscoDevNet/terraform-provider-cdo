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

func TestDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Should delete a user", func(t *testing.T) {
		httpmock.Reset()
		uid := "sample-user-uid"
		httpmock.RegisterResponder(
			"DELETE",
			"/anubis/rest/v1/users/"+uid,
			httpmock.NewJsonResponderOrPanic(200, nil),
		)
		deleteOutput, err := user.Delete(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), user.DeleteUserInput{
			Uid: uid,
		})
		assert.NotNil(t, deleteOutput, "Delete output should not be nil")
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("Should error if deletion of a user fails", func(t *testing.T) {
		httpmock.Reset()
		uid := "sample-user-uid"
		httpmock.RegisterResponder(
			"DELETE",
			"/anubis/rest/v1/users/"+uid,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		deleteOutput, err := user.Delete(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), user.DeleteUserInput{
			Uid: uid,
		})
		assert.Nil(t, deleteOutput, "Delete output should be nil")
		assert.NotNil(t, err, "error should not be nil")
	})
}
