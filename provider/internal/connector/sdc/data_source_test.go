package sdc_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testSdc = struct {
	Name string
	Uid  string
}{
	Name: "CDO_terraform-provider-cdo-SDC-1",
	Uid:  "39784a3c-0013-4e2f-af26-219560904636",
}

const testSdcTemplate = `
data "cdo_sdc_connector" "test" {
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
					resource.TestCheckResourceAttr("data.cdo_sdc_connector.test", "name", testSdc.Name),
					resource.TestCheckResourceAttr("data.cdo_sdc_connector.test", "id", testSdc.Uid),
				),
			},
		},
	})
}
