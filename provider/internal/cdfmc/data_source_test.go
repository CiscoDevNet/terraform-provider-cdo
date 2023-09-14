package cdfmc_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testCdFmc = struct {
	Hostname        string
	Uid             string
	SoftwareVersion string
}{
	Hostname:        "terraform-provider-cdo.app.staging.cdo.cisco.com",
	Uid:             "ac12b246-ed93-4a09-bc8a-5c4708854a2f",
	SoftwareVersion: "7.3.1-build 6035",
}

const testTenantTemplate = `
data "cdo_cdfmc" "test" {}`

var testTenantConfig = acctest.MustParseTemplate(testTenantTemplate, testCdFmc)

func TestAccCdFmcDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testTenantConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_cdfmc.test", "hostname", testCdFmc.Hostname),
					resource.TestCheckResourceAttr("data.cdo_cdfmc.test", "id", testCdFmc.Uid),
					resource.TestCheckResourceAttr("data.cdo_cdfmc.test", "software_version", testCdFmc.SoftwareVersion),
				),
			},
		},
	})
}
