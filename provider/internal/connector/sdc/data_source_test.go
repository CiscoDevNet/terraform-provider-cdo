package sdc_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testSdc = struct {
	TenantId string
}{
	TenantId: "ae98d25f-1089-4286-a3c5-505dcb4431a2",
}

const testSdcConfig = `
data "cdo_sdc_connector" "test" {}`

func TestAccSdcDeviceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testSdcConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_sdc_connector.test", "id", testSdc.TenantId),
				),
			},
		},
	})
}
