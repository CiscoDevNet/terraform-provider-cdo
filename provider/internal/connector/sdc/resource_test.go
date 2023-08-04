package sdc_test

import (
	"testing"

	"github.com/cisco-lockhart/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type testSdcResourceType struct {
	Name string
}

const testSdcResourceTemplate = `
resource "cdo_sdc" "test" {
	name = "{{.Name}}"
}`

var testSdcResource = testSdcResourceType{
	Name: "test-sdc-1",
}
var testSdcResourceConfig = acctest.MustParseTemplate(testSdcResourceTemplate, testSdcResource)

var testSdcResource_NewName = acctest.MustOverrideFields(testSdcResource, map[string]any{
	"Name": "test-sdc-2",
})
var testSdcResourceConfig_NewName = acctest.MustParseTemplate(testSdcResourceTemplate, testSdcResource_NewName)

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
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testSdcResourceConfig_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_sdc.test", "name", testSdcResource_NewName.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
