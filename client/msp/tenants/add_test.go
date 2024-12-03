package tenants_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/tenants"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	netHttp "net/http"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully add tenant using API token", func(t *testing.T) {
		httpmock.Reset()
		apiToken := "fake-jwt-token"
		addInp := tenants.MspAddExistingTenantInput{
			ApiToken: apiToken,
		}
		expectedResponse := tenants.MspManagedTenantStatusInfo{
			Status: "OK",
			MspManagedTenant: tenants.MspTenantOutput{
				Uid:         uuid.New().String(),
				Name:        "example-name",
				DisplayName: "Human readable name",
				Region:      "STAGING",
			},
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/api/rest/v1/msp/tenants",
			httpmock.NewJsonResponderOrPanic(201, expectedResponse),
		)

		actual, err := tenants.AddExistingTenantUsingApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), addInp)

		assert.NotNil(t, actual, "Response for added tenant should have not been nil")
		assert.Nil(t, err, "Add tenant operation should have not been an error")
		assert.Equal(t, expectedResponse.MspManagedTenant, *actual, "Add tenant operation should have returned the value of the added tenant")
	})

	t.Run("fail to add tenant using API token", func(t *testing.T) {
		httpmock.Reset()
		apiToken := "fake-jwt-token"
		addInp := tenants.MspAddExistingTenantInput{
			ApiToken: apiToken,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			"/api/rest/v1/msp/tenants",
			httpmock.NewJsonResponderOrPanic(400, nil),
		)

		actual, err := tenants.AddExistingTenantUsingApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), addInp)
		assert.Nil(t, actual, "Response for added tenant should have been nil")
		assert.NotNil(t, err, "Add tenant operation should have been an error")
	})

}
