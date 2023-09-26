package ios_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"
	"strconv"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var asaDataSourceTags = acctest.Env.IosDataSourceTags()

var testIosDataSource = struct {
	Name              string
	IgnoreCertificate string
	Tags              string
}{
	Name:              acctest.Env.IosDataSourceName(),
	IgnoreCertificate: acctest.Env.IosDataSourceIgnoreCertificate(),
	Tags:              asaDataSourceTags.AsJsonArrayString(),
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
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "tags.#", strconv.Itoa(len(asaDataSourceTags.Labels))),
					resource.TestCheckResourceAttrWith("data.cdo_ios_device.test", "tags.0", testutil.CheckEqual(asaDataSourceTags.Labels[0])),
					resource.TestCheckResourceAttrWith("data.cdo_ios_device.test", "tags.1", testutil.CheckEqual(asaDataSourceTags.Labels[1])),
					resource.TestCheckResourceAttrWith("data.cdo_ios_device.test", "tags.2", testutil.CheckEqual(asaDataSourceTags.Labels[2])),
				),
			},
		},
	})
}
