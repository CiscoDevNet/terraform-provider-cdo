package asa_test

import (
	"fmt"
	"testing"

	"github.com/cisco-lockhart/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testDeviceName      = "test-asa-device-1"
	testDeviceName2     = "test-asa-device-2"
	testIpv4            = "vasa-gb-ravpn-03-mgmt.dev.lockhart.io:443"
	testHost            = "vasa-gb-ravpn-03-mgmt.dev.lockhart.io"
	testPort            = "443"
	testSdcTypeSDC      = "SDC"
	testSdcTypeCDG      = "CDG"
	testSdcName         = "CDO_terraform-provider-cdo-SDC-1"
	testCdgUid          = "" // cdg can be figured out automatically
	testUsername        = "lockhart"
	testPassword        = "BlueSkittles123!!"
	testNewPassword     = "WrongPassword"
	testIgnoreCert      = true
	testDoNotIgnoreCert = false

	asaResourceTemplate = `resource "cdo_asa_device" "test" {
	name = %[1]q
	ipv4 = %[2]q
	sdc_type = %[3]q
	username = %[4]q
	password = %[5]q
	ignore_certificate = %[6]t
	sdc_name = %[7]q
}
`
)

var accTestAsaDeviceResourceConfig_SDC = fmt.Sprintf(asaResourceTemplate, testDeviceName, testIpv4, testSdcTypeSDC, testUsername, testPassword, testIgnoreCert, testSdcName)
var accTestAsaDeviceResourceConfig_SDC_NewName = fmt.Sprintf(asaResourceTemplate, testDeviceName2, testIpv4, testSdcTypeSDC, testUsername, testPassword, testIgnoreCert, testSdcName)
var accTestAsaDeviceResourceConfig_SDC_NewCreds = fmt.Sprintf(asaResourceTemplate, testDeviceName, testIpv4, testSdcTypeSDC, testUsername, testNewPassword, testIgnoreCert, testSdcName)

var accTestAsaDeviceResourceConfig_CDG = fmt.Sprintf(asaResourceTemplate, testDeviceName, testIpv4, testSdcTypeCDG, testUsername, testPassword, testDoNotIgnoreCert, testCdgUid)
var accTestAsaDeviceResourceConfig_CDG_NewName = fmt.Sprintf(asaResourceTemplate, testDeviceName2, testIpv4, testSdcTypeCDG, testUsername, testPassword, testDoNotIgnoreCert, testCdgUid)
var accTestAsaDeviceResourceConfig_CDG_NewCreds = fmt.Sprintf(asaResourceTemplate, testDeviceName, testIpv4, testSdcTypeCDG, testUsername, testNewPassword, testDoNotIgnoreCert, testCdgUid)

func TestAccAsaDeviceResource_SDC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + accTestAsaDeviceResourceConfig_SDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testDeviceName),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "ipv4", testIpv4),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "host", testHost),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "port", testPort),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "sdc_type", testSdcTypeSDC),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "username", testUsername),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testPassword),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + accTestAsaDeviceResourceConfig_SDC_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testDeviceName2),
				),
			},
			{
				Config: acctest.ProviderConfig() + accTestAsaDeviceResourceConfig_SDC_NewCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testNewPassword),
				),
			},
			{
				Config: acctest.ProviderConfig() + accTestAsaDeviceResourceConfig_SDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testPassword),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAsaDeviceResource_CDG(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + accTestAsaDeviceResourceConfig_CDG,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testDeviceName),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "ipv4", testIpv4),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "host", testHost),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "port", testPort),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "sdc_type", testSdcTypeCDG),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "username", testUsername),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testPassword),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + accTestAsaDeviceResourceConfig_CDG_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testDeviceName2),
				),
			},
			{
				Config: acctest.ProviderConfig() + accTestAsaDeviceResourceConfig_CDG_NewCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testNewPassword),
				),
			},
			{
				Config: acctest.ProviderConfig() + accTestAsaDeviceResourceConfig_CDG,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testPassword),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
