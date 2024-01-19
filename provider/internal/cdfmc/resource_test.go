package cdfmc_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var resourceModel = struct {
	Name       string
	Hostname   string
	DomainUuid string
}{
	Name:     acctest.Env.CloudFmcResourceName(),
	Hostname: acctest.Env.CloudFmcResourceHostname(),
}

const resourceTemplate = `
resource "cdo_cdfmc" "test" {
}`

var resourceConfig = acctest.MustParseTemplate(resourceTemplate, resourceModel)

func TestAccCdFmcResource(t *testing.T) {
	t.Skip("we cant delete a fmc so this test cannot be run, or we should find a way to spin up new environment, either seems uneasy, skipping for now.")
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_cdfmc.test", "name", resourceModel.Name),
					resource.TestCheckResourceAttr("cdo_cdfmc.test", "hostname", resourceModel.Hostname),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
