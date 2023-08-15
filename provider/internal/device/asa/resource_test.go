package asa_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type testAsaResourceType struct {
	Name              string
	Ipv4              string
	SdcType           string
	Username          string
	Password          string
	IgnoreCertificate bool
	SdcName           string

	Host string
	Port string
}

const asaResourceTemplate = `
resource "cdo_asa_device" "test" {
	name = "{{.Name}}"
	socket_address = "{{.SocketAddress}}"
	connector_type = "{{.SdcType}}"
	username = "{{.EncryptedUsername}}"
	password = "{{.EncryptedPassword}}"
	ignore_certificate = "{{.IgnoreCertificate}}"
	sdc_name = "{{.SdcName}}"
}`

// SDC configs.

// default config.
var testAsaResource_SDC = testAsaResourceType{
	Name:              "test-asa-device-1",
	Ipv4:              "vasa-gb-ravpn-03-mgmt.dev.lockhart.io:443",
	SdcType:           "SDC",
	Username:          "lockhart",
	Password:          "BlueSkittles123!!",
	IgnoreCertificate: true,
	SdcName:           "CDO_terraform-provider-cdo-SDC-1",

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
var testAsaResource_SDC_NewCreds = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{
	"EncryptedPassword": "WrongPassword",
})
var testAsaResourceConfig_SDC_NewCreds = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_NewCreds)

var testAsaResource_SDC_NewLocation = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{"SocketAddress": alternativeDeviceLocation})
var testAsaResourceConfig_SDC_NewLocation = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_NewLocation)

// CDG configs

// default config.
var testAsaResource_CDG = testAsaResourceType{
	Name:              "test-asa-device-1",
	Ipv4:              "vasa-gb-ravpn-03-mgmt.dev.lockhart.io:443",
	SdcType:           "CDG",
	Username:          "lockhart",
	Password:          "BlueSkittles123!!",
	IgnoreCertificate: false,
	SdcName:           "CDO_terraform-provider-cdo-SDC-1",

	Host: "vasa-gb-ravpn-03-mgmt.dev.lockhart.io",
	Port: "443",
}
var testAsaResourceConfig_CDG = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG)

// new name config.
var testAsaResource_CDG_NewName = acctest.MustOverrideFields(testAsaResource_CDG, map[string]any{
	"Name": "test-asa-device-2",
})
var testAsaResourceConfig_CDG_NewName = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG_NewName)

// new creds config.
var testAsaResource_CDG_NewCreds = acctest.MustOverrideFields(testAsaResource_CDG, map[string]any{
	"EncryptedPassword": "WrongPassword",
})
var testAsaResourceConfig_CDG_NewCreds = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG_NewCreds)

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
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_SDC.Ipv4),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "host", testAsaResource_SDC.Host),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "port", testAsaResource_SDC.Port),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "connector_type", testAsaResource_SDC.SdcType),
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
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC_NewCreds.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_SDC_NewCreds.Password),
				),
			},
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_SDC.Password),
				),
			},
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewLocation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "ipv4", testAsaResource_SDC_NewLocation.Ipv4),
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
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_CDG.Ipv4),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "host", testAsaResource_CDG.Host),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "port", testAsaResource_CDG.Port),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "connector_type", testAsaResource_CDG.SdcType),
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
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG_NewCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG_NewCreds.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_CDG_NewCreds.Password),
				),
			},
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_CDG.Password),
				),
			},
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG_NewLocation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "ipv4", testAsaResource_CDG_NewLocation.Ipv4),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
