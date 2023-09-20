package ios_test

import (
	"strconv"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testIosDataSource = struct {
	Id                string
	Name              string
	SocketAddress     string
	Host              string
	Port              int64
	IgnoreCertificate string
}{
	Id:                acctest.Env.IosDataSourceId(),
	Name:              acctest.Env.IosDataSourceName(),
	SocketAddress:     acctest.Env.IosDataSourceSocketAddress(),
	Host:              acctest.Env.IosDataSourceHost(),
	Port:              acctest.Env.IosDataSourcePort(),
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
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "id", testIosDataSource.Id),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "name", testIosDataSource.Name),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "socket_address", testIosDataSource.SocketAddress),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "host", testIosDataSource.Host),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "port", strconv.FormatInt(testIosDataSource.Port, 10)),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "ignore_certificate", testIosDataSource.IgnoreCertificate),
				),
			},
		},
	})
}
