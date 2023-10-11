package examples_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var resourceModel = struct {
	Name string
}{
	Name: "example name", // acctest.Env.ExampleResourceName(),
}

const resourceTemplate = `
resource "cdo_example" "test" {
	name = "{{.Name}}"
}`

var resourceConfig = acctest.MustParseTemplate(resourceTemplate, resourceModel)

var resourceModel_NewName = acctest.MustOverrideFields(resourceModel, map[string]any{
	"Name": "example new name", // acctest.Env.ExampleResourceNewName(),
})
var resourceConfig_NewName = acctest.MustParseTemplate(resourceTemplate, resourceModel_NewName)

func TestAccExampleResource(t *testing.T) {
	t.Skip("this is an example resource test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_example.test", "name", resourceModel.Name),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + resourceConfig_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_example.test", "name", resourceModel_NewName.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
