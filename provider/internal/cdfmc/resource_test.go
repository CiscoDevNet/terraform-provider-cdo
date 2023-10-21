package cdfmc_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var resourceModel = struct {
	Name string
}{
	Name: acctest.Env.CloudFmcResourceName(),
}

const resourceTemplate = `
resource "cdo_cdfmc" "test" {
}`

var resourceConfig = acctest.MustParseTemplate(resourceTemplate, resourceModel)

func TestAccCdFmcResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_cdfmc.test", "name", resourceModel.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
