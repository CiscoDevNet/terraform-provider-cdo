package cdfmc_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testCdFmc = struct {
	Hostname string
}{
	Hostname: acctest.Env.CloudFmcDataSourceHostname(),
}

const testCdFmcTemplate = `
data "cdo_cdfmc" "test" {}`

var testCdfmcConfig = acctest.MustParseTemplate(testCdFmcTemplate, testCdFmc)

func TestAccCdFmcDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testCdfmcConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_cdfmc.test", "hostname", testCdFmc.Hostname),
				),
			},
		},
	})
}
