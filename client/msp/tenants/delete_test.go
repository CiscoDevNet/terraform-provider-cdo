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

func TestDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully delete tenant", func(t *testing.T) {
		httpmock.Reset()
		var tenantUid = uuid.New().String()
		var deleteInput = tenants.DeleteByUidInput{
			Uid: tenantUid,
		}
		httpmock.RegisterResponder(
			netHttp.MethodDelete,
			"/api/rest/v1/msp/tenants/"+tenantUid,
			httpmock.NewJsonResponderOrPanic(204, nil),
		)
		response, err := tenants.DeleteByUid(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), deleteInput)
		assert.Nil(t, response)
		assert.Nil(t, err)
	})

	t.Run("send through error if deletion failed", func(t *testing.T) {
		httpmock.Reset()
		var tenantUid = uuid.New().String()
		var deleteInput = tenants.DeleteByUidInput{
			Uid: tenantUid,
		}
		httpmock.RegisterResponder(
			netHttp.MethodDelete,
			"/api/rest/v1/msp/tenants/"+tenantUid,
			httpmock.NewJsonResponderOrPanic(500, "Not found"),
		)
		response, err := tenants.DeleteByUid(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), deleteInput)
		assert.Nil(t, response)
		assert.NotNil(t, err)
	})
}
