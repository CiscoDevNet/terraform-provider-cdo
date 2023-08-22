package ios_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type testIosResourceType struct {
	Name              string
	SocketAddress     string
	ConnectorType     string
	Username          string
	Password          string
	ConnectorName     string
	IgnoreCertificate bool

	Host string
	Port string
}

const testIosResourceTemplate = `
resource "cdo_ios_device" "test" {
	name = "{{.Name}}"
	socket_address = "{{.SocketAddress}}"
	username = "{{.Username}}"
	password = "{{.Password}}"
	connector_name = "{{.ConnectorName}}"
	ignore_certificate = "{{.IgnoreCertificate}}"
}`

var testIosResource = testIosResourceType{
	Name:              "test-ios-device-1",
	SocketAddress:     "10.10.0.198:22",
	ConnectorType:     "SDC",
	Username:          "lockhart",
	Password:          "BlueSkittles123!!",
	ConnectorName:     "CDO_terraform-provider-cdo-SDC-1",
	IgnoreCertificate: true,

	Host: "10.10.0.198",
	Port: "22",
}
var testIosResourceConfig = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource)

var testIosResource_NewName = acctest.MustOverrideFields(testIosResource, map[string]any{
	"Name": "test-ios-device-2",
})
var testIosResourceConfig_NewName = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource_NewName)

func TestAccIosDeviceResource_SDC(t *testing.T) {
	t.Skip("require new ios device in ci")
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testIosResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testIosResource.Name),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "socket_address", testIosResource.SocketAddress),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "host", testIosResource.Host),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "port", testIosResource.Port),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "connector_type", testIosResource.ConnectorType),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "username", testIosResource.Username),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "password", testIosResource.Password),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testIosResourceConfig_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testIosResource_NewName.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
