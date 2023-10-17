package ftd_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var dataSourceModel = struct {
	Name             string
	AccessPolicyName string
	PerformanceTier  string
	Virtual          string
	Licenses         string
}{
	Name:             acctest.Env.FtdDataSourceName(), // typically we get the actual value from the environment like this: acctest.Env.ExampleDataSourceName(), so that most parameters are configurable so that it can be run in different CDO environment
	AccessPolicyName: acctest.Env.FtdDataSourceAccessPolicyName(),
	PerformanceTier:  acctest.Env.FtdDataSourcePerformanceTier(),
	Virtual:          acctest.Env.FtdDataSourceVirtual(),
	Licenses:         acctest.Env.FtdDataSourceLicenses(),
}
var dataSourceTemplate = `
data "cdo_ftd_device" "test" {
	name = "{{.Name}}"
}`
var dataSourceConfig = acctest.MustParseTemplate(dataSourceTemplate, dataSourceModel)

func TestAccFtdDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_ftd_device.test", "name", dataSourceModel.Name),
					resource.TestCheckResourceAttr("data.cdo_ftd_device.test", "access_policy_name", dataSourceModel.AccessPolicyName),
					resource.TestCheckResourceAttr("data.cdo_ftd_device.test", "performance_tier", dataSourceModel.PerformanceTier),
					resource.TestCheckResourceAttr("data.cdo_ftd_device.test", "virtual", dataSourceModel.Virtual),
					resource.TestCheckResourceAttr("data.cdo_ftd_device.test", "licenses.#", "1"), // number of licenses = 1
				),
			},
		},
	})
}
