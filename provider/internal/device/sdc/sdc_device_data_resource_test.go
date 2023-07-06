package sdc_test

import (
	"testing"

	"github.com/cisco-lockhart/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	// TODO: refactor to some constants.go file before things becomes messy.
	accTestTenantId = "ae98d25f-1089-4286-a3c5-505dcb4431a2" // this is real sdc device id in ci
)

func TestAccSdcDeviceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + accTestSdcDeviceDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_sdc_device.test", "id", accTestTenantId),
				),
			},
		},
	})
}

const accTestSdcDeviceDataSourceConfig = `
data "cdo_sdc_device" "test" {}
`
