package asa_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAsaDataSource = struct {
	Id                string
	SdcType           string
	Name              string
	Ipv4              string
	Host              string
	Port              string
	IgnoreCertificate string
}{
	Id:                "331ff184-9ae6-45f3-8c55-71a150a6b58f",
	SdcType:           "CDG",
	Name:              "asa-data-source",
	Ipv4:              "52.53.230.145:443",
	Host:              "52.53.230.145",
	Port:              "443",
	IgnoreCertificate: "false",
}

const testAsaDataSourceTemplate = `
data "cdo_asa_device" "test" {
	id = "{{.Id}}"
}`

var testAsaDataSourceConfig = acctest.MustParseTemplate(testAsaDataSourceTemplate, testAsaDataSource)

func TestAccAsaDeviceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testAsaDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "id", testAsaDataSource.Id),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "connector_type", testAsaDataSource.SdcType),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "name", testAsaDataSource.Name),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "socket_address", testAsaDataSource.Ipv4),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "host", testAsaDataSource.Host),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "port", testAsaDataSource.Port),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "ignore_certificate", testAsaDataSource.IgnoreCertificate),
				),
			},
		},
	})
}
