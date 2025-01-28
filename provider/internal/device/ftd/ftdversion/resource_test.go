package ftdversion_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

var dataSourceModel = struct {
	Name string
}{
	Name: acctest.Env.FtdDataSourceName(), // typically we get the actual value from the environment like this: acctest.Env.ExampleDataSourceName(), so that most parameters are configurable so that it can be run in different CDO environment
}

var sameVersionTemplate = `
data "cdo_ftd_device" "test" {
	name = "{{.Name}}"
}

resource "cdo_ftd_device_version" "test" {
	ftd_uid = data.cdo_ftd_device.test.id
	software_version = "7.3.0"
}
`

var config = acctest.MustParseTemplate(sameVersionTemplate, dataSourceModel)

func TestAccFtdVersionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// when the software_version specified in cdo_ftd_device_version is the same as the version on the FTD, then I should not fail
			{
				Config: acctest.ProviderConfig() + config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ftd_device_version.test", "software_version_on_device", "7.3.0"),
				),
			},
		},
	})
}
