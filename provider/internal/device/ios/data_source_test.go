package ios_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testIosDataSource = struct {
	Id                string
	ConnectorType     string
	Name              string
	SocketAddress     string
	Host              string
	Port              string
	IgnoreCertificate string
}{
	Id:                "cd0483d0-5ec5-4d8e-b92d-8eb389f88417",
	ConnectorType:     "SDC",
	Name:              "weilue-test-ios",
	SocketAddress:     "10.10.0.198:22",
	Host:              "10.10.0.198",
	Port:              "22",
	IgnoreCertificate: "false",
}

var testIosDataSourceTemplate = `
data "cdo_ios_device" "test" {
	id = "{{.Id}}"
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
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "connector_type", testIosDataSource.ConnectorType),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "name", testIosDataSource.Name),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "socket_address", testIosDataSource.SocketAddress),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "host", testIosDataSource.Host),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "port", testIosDataSource.Port),
					resource.TestCheckResourceAttr("data.cdo_ios_device.test", "ignore_certificate", testIosDataSource.IgnoreCertificate),
				),
			},
		},
	})
}
