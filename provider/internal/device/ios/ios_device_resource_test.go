package ios_test

import (
	"fmt"
	"testing"

	"github.com/cisco-lockhart/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testDeviceName      = "test-ios-device-1"
	testDeviceName2     = "test-ios-device-2"
	testIpv4            = "10.10.0.198:22"
	testHost            = "10.10.0.198"
	testPort            = "22"
	testSdcTypeSDC      = "SDC"
	testSdcTypeCDG      = "CDG"
	testSdcUid          = "39784a3c-0013-4e2f-af26-219560904636"
	testCdgUid          = "" // cdg can be figured out automatically
	testUsername        = "lockhart"
	testPassword        = "BlueSkittles123!!"
	testIgnoreCert      = true
	testDoNotIgnoreCert = false

	iosResourceTemplate = `resource "cdo_ios_device" "test" {
	name = %[1]q
	ipv4 = %[2]q
	sdc_type = %[3]q
	username = %[4]q
	password = %[5]q
	ignore_certificate = %[6]t
	sdc_uid = %[7]q
}
`
)

var accTestIosDeviceResourceConfig_SDC = fmt.Sprintf(iosResourceTemplate, testDeviceName, testIpv4, testSdcTypeSDC, testUsername, testPassword, testIgnoreCert, testSdcUid)
var accTestIosDeviceResourceConfig_SDC_NewName = fmt.Sprintf(iosResourceTemplate, testDeviceName2, testIpv4, testSdcTypeSDC, testUsername, testPassword, testIgnoreCert, testSdcUid)

var accTestIosDeviceResourceConfig_CDG = fmt.Sprintf(iosResourceTemplate, testDeviceName, testIpv4, testSdcTypeCDG, testUsername, testPassword, testDoNotIgnoreCert, testCdgUid)
var accTestIosDeviceResourceConfig_CDG_NewName = fmt.Sprintf(iosResourceTemplate, testDeviceName2, testIpv4, testSdcTypeCDG, testUsername, testPassword, testDoNotIgnoreCert, testCdgUid)

func TestAccIosDeviceResource_SDC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + accTestIosDeviceResourceConfig_SDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testDeviceName),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "ipv4", testIpv4),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "host", testHost),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "port", fmt.Sprint(testPort)),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "sdc_type", testSdcTypeSDC),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "username", testUsername),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "password", testPassword),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + accTestIosDeviceResourceConfig_SDC_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testDeviceName2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIosDeviceResource_CDG(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + accTestIosDeviceResourceConfig_CDG,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testDeviceName),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "ipv4", testIpv4),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "host", testHost),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "port", fmt.Sprint(testPort)),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "sdc_type", testSdcTypeCDG),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "username", testUsername),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "password", testPassword),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + accTestIosDeviceResourceConfig_CDG_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testDeviceName2),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
