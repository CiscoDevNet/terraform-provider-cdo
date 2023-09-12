package tenant_test

import (
	"context"
	netHttp "net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/tenant"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestReadTenantDetails(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Should get current tenant details", func(t *testing.T) {
		httpmock.Reset()
		expected := tenant.ReadTenantDetailsOutput{
			UserAuthentication: tenant.UserAuthentication{
				Details: tenant.TenantDetailsDetails{
					TenantUid:              "111-111-111-111",
					TenantName:             "sample-tenant-name",
					TenantOrganizationName: "sample-org-name",
					TenantPayType:          "NOT_PAYING",
				},
			},
		}
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/oauth/check_token",
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		actual, err := tenant.ReadTenantDetails(context.Background(), *http.MustNewWithConfig(baseUrl, "valid token", 0, 0, time.Minute))
		assert.NotNil(t, actual, "Read output should not be nil")
		assert.Equal(t, *actual, expected)
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("should error if getting current tenant details fails", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/anubis/rest/v1/oauth/check_token",
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := tenant.ReadTenantDetails(context.Background(), *http.MustNewWithConfig(baseUrl, "valid token", 0, 0, time.Minute))
		assert.Nil(t, actual, "Read output should be nil")
		assert.NotNil(t, err, "error should not be nil")
	})
}
