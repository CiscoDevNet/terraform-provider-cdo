package ios_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testIosDataSource = struct {
	Name              string
	IgnoreCertificate string
}{
	Name:              acctest.Env.IosDataSourceName(),
	IgnoreCertificate: acctest.Env.IosDataSourceIgnoreCertificate(),
}

var testIosDataSourceTemplate = `
data "cdo_ios_device" "test" {
	name = "{{.Name}}"
}`
var testIosDataSourceConfig = acctest.MustParseTemplate(testIosDataSourceTemplate, testIosDataSource)

func TestAccIosDeviceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testIosDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "name", testIosDataSource.Name),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "ignore_certificate", testIosDataSource.IgnoreCertificate),
				),
			},
		},
	})
}
