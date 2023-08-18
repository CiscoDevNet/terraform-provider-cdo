package asa_test

import (
	"regexp"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type testAsaResourceType struct {
	Name              string
	SocketAddress     string
	ConnectorName     string
	ConnectorType     string
	Username          string
	Password          string
	IgnoreCertificate bool

	Host string
	Port string
}

const asaResourceTemplate = `
resource "cdo_asa_device" "test" {
	name = "{{.Name}}"
	socket_address = "{{.SocketAddress}}"
	connector_name = "{{.ConnectorName}}"
	connector_type = "{{.ConnectorType}}"
	username = "{{.Username}}"
	password = "{{.Password}}"
	ignore_certificate = "{{.IgnoreCertificate}}"
}`

// SDC configs.

// default config.
var testAsaResource_SDC = testAsaResourceType{
	Name:              "test-asa-device-1",
	SocketAddress:     "vasa-gb-ravpn-03-mgmt.dev.lockhart.io:443",
	ConnectorName:     "CDO_terraform-provider-cdo-SDC-1",
	ConnectorType:     "SDC",
	Username:          "lockhart",
	Password:          "BlueSkittles123!!",
	IgnoreCertificate: true,

	Host: "vasa-gb-ravpn-03-mgmt.dev.lockhart.io",
	Port: "443",
}

const alternativeDeviceLocation = "35.177.20.218:443"

var testAsaResourceConfig_SDC = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC)

// new name config.
var testAsaResource_SDC_NewName = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{
	"Name": "test-asa-device-2",
})
var testAsaResourceConfig_SDC_NewName = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_NewName)

// new creds config.
var testAsaResource_SDC_BadCreds = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{
	"Password": "WrongPassword",
})
var testAsaResourceConfig_SDC_NewCreds = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_BadCreds)

var testAsaResource_SDC_NewLocation = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{"SocketAddress": alternativeDeviceLocation})
var testAsaResourceConfig_SDC_NewLocation = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_NewLocation)

// CDG configs

// default config.
var testAsaResource_CDG = testAsaResourceType{
	Name:              "test-asa-device-1",
	SocketAddress:     "vasa-gb-ravpn-03-mgmt.dev.lockhart.io:443",
	ConnectorName:     "CDO_terraform-provider-cdo-SDC-1",
	ConnectorType:     "CDG",
	Username:          "lockhart",
	Password:          "BlueSkittles123!!",
	IgnoreCertificate: false,

	Host: "vasa-gb-ravpn-03-mgmt.dev.lockhart.io",
	Port: "443",
}
var testAsaResourceConfig_CDG = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG)

// new name config.
var testAsaResource_CDG_NewName = acctest.MustOverrideFields(testAsaResource_CDG, map[string]any{
	"Name": "test-asa-device-2",
})
var testAsaResourceConfig_CDG_NewName = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG_NewName)

// bad credentials config.
var testAsaResource_CDG_BadCreds = acctest.MustOverrideFields(testAsaResource_CDG, map[string]any{
	"Password": "WrongPassword",
})
var testAsaResourceConfig_CDG_NewCreds = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG_BadCreds)

var testAsaResource_CDG_NewLocation = acctest.MustOverrideFields(testAsaResource_CDG, map[string]any{"SocketAddress": alternativeDeviceLocation})
var testAsaResourceConfig_CDG_NewLocation = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG_NewLocation)

func TestAccAsaDeviceResource_SDC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_SDC.SocketAddress),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "host", testAsaResource_SDC.Host),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "port", testAsaResource_SDC.Port),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "connector_type", testAsaResource_SDC.ConnectorType),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "username", testAsaResource_SDC.Username),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_SDC.Password),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC_NewName.Name),
				),
			},
			// bad credential tests
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC_BadCreds.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_SDC_BadCreds.Password),
				),
				ExpectError: regexp.MustCompile(`.*bad credentials.*`),
			},
			// fix correct credentials test
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_SDC.Password),
				),
			},
			// change location test
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewLocation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_SDC_NewLocation.SocketAddress),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAsaDeviceResource_CDG(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_CDG.SocketAddress),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "host", testAsaResource_CDG.Host),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "port", testAsaResource_CDG.Port),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "connector_type", testAsaResource_CDG.ConnectorType),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "username", testAsaResource_CDG.Username),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_CDG.Password),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG_NewName.Name),
				),
			},
			// bad credentials tests
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG_NewCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG_BadCreds.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_CDG_BadCreds.Password),
				),
				ExpectError: regexp.MustCompile(".*bad credentials.*"),
			},
			// fix bad credentials test
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_CDG.Password),
				),
			},
			// change location test
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG_NewLocation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_CDG_NewLocation.SocketAddress),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
