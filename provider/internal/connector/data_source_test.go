package connector_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testSdc = struct {
	Name string
}{
	Name: acctest.Env.ConnectorDataSourceName(),
}

const testSdcTemplate = `
data "cdo_sdc" "test" {
	name = "{{.Name}}"
}`

var testSdcConfig = acctest.MustParseTemplate(testSdcTemplate, testSdc)

func TestAccSdcDeviceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testSdcConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_sdc.test", "name", testSdc.Name),
				),
			},
		},
	})
}
