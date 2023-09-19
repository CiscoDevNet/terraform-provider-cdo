package connector_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type testSdcResourceType struct {
	Name string
}

const testResourceTemplate = `
resource "cdo_sdc" "test" {
	name = "{{.Name}}"
}`

var testSdcResource = testSdcResourceType{
	Name: acctest.Env.ConnectorResourceName(),
}
var testSdcResourceConfig = acctest.MustParseTemplate(testResourceTemplate, testSdcResource)

var testResource_NewName = acctest.MustOverrideFields(testSdcResource, map[string]any{
	"Name": acctest.Env.ConnectorResourceNewName(),
})
var testResourceConfig_NewName = acctest.MustParseTemplate(testResourceTemplate, testResource_NewName)

func TestAccSdcResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testSdcResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_sdc.test", "name", testSdcResource.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccUpdateSdcResource(t *testing.T) {
	t.Skip("Requires us to figure out how to wait for the resource to finish updating")
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testSdcResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_sdc.test", "name", testSdcResource.Name),
				),
			},
			// Update and Read testing
			// commenting out because this test is flaking in CI
			{
				Config: acctest.ProviderConfig() + testResourceConfig_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_sdc.test", "name", testResource_NewName.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
