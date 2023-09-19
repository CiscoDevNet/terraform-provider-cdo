package user_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type testUserResourceType struct {
	Name        string
	ApiOnlyUser string
	UserRole    string
}

const testResourceTemplate = `
resource "cdo_user" "test" {
	name = "{{.Name}}"
    is_api_only_user = "{{.ApiOnlyUser}}"
    role = "{{.UserRole}}"
}`

var testUserResource = testUserResourceType{
	Name:        acctest.Env.UserResourceName(),
	ApiOnlyUser: acctest.Env.UserResourceIsApiOnly(),
	UserRole:    acctest.Env.UserResourceRole(),
}
var testUserResourceConfig = acctest.MustParseTemplate(testResourceTemplate, testUserResource)
var testResource_NewName = acctest.MustOverrideFields(testUserResource, map[string]any{
	"Name": acctest.Env.UserResourceNewName(),
})
var testResourceConfig_NewName = acctest.MustParseTemplate(testResourceTemplate, testResource_NewName)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testUserResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_user.test", "name", testUserResource.Name),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testResourceConfig_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_user.test", "name", testResource_NewName.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
