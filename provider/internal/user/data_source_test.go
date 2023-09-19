package user_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testUserName = acctest.Env.UserDataSourceUser()
var testRole = acctest.Env.UserDataSourceRole()
var testIsApiOnlyUser = acctest.Env.UserDataSourceIsApiOnly()

var testUser = struct {
	Name        string
	ApiOnlyUser string
	UserRole    string
}{
	Name:        testUserName,      //"terraform-provider-cdo-super-admin@lockhart.io",
	ApiOnlyUser: testIsApiOnlyUser, // "false"
	UserRole:    testRole,          //"ROLE_SUPER_ADMIN",
}

const testUserTemplate = `
data "cdo_user" "test" {
	name = "{{.Name}}"
}`

var testUserConfig = acctest.MustParseTemplate(testUserTemplate, testUser)

func TestAccUserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testUserConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_user.test", "name", testUser.Name),
					resource.TestCheckResourceAttr("data.cdo_user.test", "is_api_only_user", testUser.ApiOnlyUser),
					resource.TestCheckResourceAttr("data.cdo_user.test", "role", testUser.UserRole),
				),
			},
		},
	})
}
