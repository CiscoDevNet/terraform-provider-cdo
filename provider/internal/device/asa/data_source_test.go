package asa_test

import (
	"strconv"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAsaDataSource = struct {
	ConnectorType     string
	Name              string
	SocketAddress     string
	Host              string
	Port              int64
	IgnoreCertificate bool
	Tags              []string
}{
	ConnectorType:     acctest.Env.AsaDataSourceConnectorType(),
	Name:              acctest.Env.AsaDataSourceName(),
	SocketAddress:     acctest.Env.AsaDataSourceSocketAddress(),
	Host:              acctest.Env.AsaDataSourceHost(),
	Port:              acctest.Env.AsaDataSourcePort(),
	IgnoreCertificate: acctest.Env.AsaDataSourceIgnoreCertificate(),
	Tags:              acctest.Env.AsaDataSourceTags(),
}

const testAsaDataSourceTemplate = `
data "cdo_asa_device" "test" {
	name = "{{.Name}}"
}`

var testAsaDataSourceConfig = acctest.MustParseTemplate(testAsaDataSourceTemplate, testAsaDataSource)

func TestAccAsaDeviceDataSource(t *testing.T) {
    t.Skip("LH-84085")
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testAsaDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "connector_type", testAsaDataSource.ConnectorType),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "name", testAsaDataSource.Name),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "socket_address", testAsaDataSource.SocketAddress),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "host", testAsaDataSource.Host),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "port", strconv.FormatInt(testAsaDataSource.Port, 10)),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "ignore_certificate", strconv.FormatBool(testAsaDataSource.IgnoreCertificate)),
					resource.TestCheckResourceAttr("data.cdo_asa_device.test", "labels.#", strconv.Itoa(len(testAsaDataSource.Tags))),
					resource.TestCheckResourceAttrWith("data.cdo_asa_device.test", "labels.0", testutil.CheckEqual(testAsaDataSource.Tags[0])),
				),
			},
		},
	})
}
