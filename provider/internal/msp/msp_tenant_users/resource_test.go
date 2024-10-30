package msp_tenant_users_test

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

type Users struct {
	Username    string
	Role        string
	ApiOnlyUser bool
}

var testMspManagedTenantUsersResource = struct {
	TenantUid string
	Users     []Users
}{
	Users: []Users{
		{Username: "user1@example.com", Role: "ROLE_SUPER_ADMIN", ApiOnlyUser: false},
		{Username: "example-api-user", Role: "ROLE_ADMIN", ApiOnlyUser: true},
	},
	TenantUid: acctest.Env.MspTenantId(),
}

const testMspManagedTenantUsersTemplate = `
resource "cdo_msp_managed_tenant_users" "test" {
	tenant_uid = "{{.TenantUid}}"
	users = [
		{
			"username": "{{(index .Users 0).Username}}"
			"role": "{{(index .Users 0).Role}}"
			"api_only_user": "{{(index .Users 0).ApiOnlyUser}}"
		},
		{
			"username": "{{(index .Users 1).Username}}"
			"role": "{{(index .Users 1).Role}}"
			"api_only_user": {{(index .Users 1).ApiOnlyUser}}
		}
	]
}`

var testMspManagedTenantUsersResourceConfig = acctest.MustParseTemplate(testMspManagedTenantUsersTemplate, testMspManagedTenantUsersResource)

func TestAccMspManagedTenantUsersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.MspProviderConfig() + testMspManagedTenantUsersResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_users.test", "tenant_uid", testMspManagedTenantUsersResource.TenantUid),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_users.test", "users.0.username", testMspManagedTenantUsersResource.Users[0].Username),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_users.test", "users.0.role", testMspManagedTenantUsersResource.Users[0].Role),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_users.test", "users.0.api_only_user", fmt.Sprintf("%t", testMspManagedTenantUsersResource.Users[0].ApiOnlyUser)),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_users.test", "users.1.username", testMspManagedTenantUsersResource.Users[1].Username),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_users.test", "users.1.role", testMspManagedTenantUsersResource.Users[1].Role),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_users.test", "users.1.api_only_user", fmt.Sprintf("%t", testMspManagedTenantUsersResource.Users[1].ApiOnlyUser)),
				),
			},
		},
	})
}
