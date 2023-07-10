package ios_test

import (
	"fmt"
	"testing"

	"github.com/cisco-lockhart/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testDeviceName  = "test-ios-device-1"
	testDeviceName2 = "test-ios-device-2"
	testIpv4        = "10.10.6.53:22"
	testHost        = "10.10.6.53"
	testPort        = "22"
	testSdcTypeSDC  = "SDC"
	testSdcUid      = "bc58aa81-0ac5-427d-beff-a3fc0e8f65c6"
	testUsername    = "lockhart"
	testPassword    = "BlueSkittles123!!"
	testIgnoreCert  = true

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

func TestAccIosDeviceResource_SDC(t *testing.T) {
	t.Skip("requires SDC set up in CI")

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
