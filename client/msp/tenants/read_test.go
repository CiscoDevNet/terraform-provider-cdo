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

func TestRead(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully read tenant", func(t *testing.T) {
		httpmock.Reset()
		var entityUid = uuid.New().String()
		var tenantResponse = tenants.MspTenantOutput{
			Uid:         entityUid,
			Name:        "test-tenant",
			DisplayName: "Pineapple Crushers Inc",
			Region:      "STAGING",
		}
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/api/rest/v1/msp/tenants/"+entityUid,
			httpmock.NewJsonResponderOrPanic(200, tenantResponse),
		)
		actual, err := tenants.Read(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), tenants.ReadByUidInput{
			Uid: entityUid,
		})

		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})

	t.Run("fail to read tenant", func(t *testing.T) {
		httpmock.Reset()
		var entityUid = uuid.New().String()
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/api/rest/v1/msp/tenants/"+entityUid,
			httpmock.NewJsonResponderOrPanic(404, "Not found"),
		)
		actual, err := tenants.Read(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), tenants.ReadByUidInput{
			Uid: entityUid,
		})

		assert.Nil(t, actual)
		assert.ErrorContains(t, err, "Not Found")
	})
}
