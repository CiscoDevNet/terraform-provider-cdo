package msp_tenant_user_api_token_test

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"strings"
	"testing"
	"text/template"
)

type Users struct {
	Username    string
	Roles       []string
	ApiOnlyUser bool
}

var testMspManagedTenantUsersResource = struct {
	TenantUid string
	Users     []Users
}{
	Users: []Users{
		{Username: "api_only_user", Roles: []string{"ROLE_SUPER_ADMIN"}, ApiOnlyUser: true},
	},
	TenantUid: acctest.Env.MspTenantId(),
}

// Join function to concatenate elements of a slice into a JSON array string.
func join(slice []string) string {
	quoted := make([]string, len(slice))
	for i, s := range slice {
		quoted[i] = fmt.Sprintf("%q", s) // Quotes each role to make it valid JSON
	}
	return strings.Join(quoted, ", ") // Joins with a comma
}

const testMspManagedTenantUsersAndApiTokenTemplate = `
resource "cdo_msp_managed_tenant_users" "test" {
	tenant_uid = "{{.TenantUid}}"
	users = [
		{
			"username": "{{(index .Users 0).Username}}"
			"roles": [{{ join (index .Users 0).Roles }}]
			"api_only_user": "{{(index .Users 0).ApiOnlyUser}}"
		}
	]
}

resource "cdo_msp_managed_tenant_user_api_token" "test" {
	tenant_uid = "{{.TenantUid}}"
	user_uid = cdo_msp_managed_tenant_users.test.users[0].id
}
`

var testMspManagedTenantUsersAndApiTokenResourceConfig = acctest.MustParseTemplateWithFuncMap(testMspManagedTenantUsersAndApiTokenTemplate, testMspManagedTenantUsersResource, template.FuncMap{
	"join": join,
})

func TestAccMspManagedTenantUserApiTokenResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.MspProviderConfig() + testMspManagedTenantUsersAndApiTokenResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cdo_msp_managed_tenant_user_api_token.test", "api_token"),
				),
			},
		},
	})
}
