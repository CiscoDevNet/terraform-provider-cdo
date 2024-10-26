package msp_tenant_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

var testMspTenantDataSource = struct {
	Name        string
	DisplayName string
	Id          string
	Region      string
}{
	Name:        acctest.Env.MspTenantName(),
	DisplayName: acctest.Env.MspTenantDisplayName(),
	Id:          acctest.Env.MspTenantId(),
	Region:      acctest.Env.MspTenantRegion(),
}

const testMspTenantDataSourceTemplate = `
data "cdo_msp_managed_tenant" "test" {
	name = "{{.Name}}"
}`

var testMspTenantDataSourceConfig = acctest.MustParseTemplate(testMspTenantDataSourceTemplate, testMspTenantDataSource)

func TestAccMspTenantDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.MspProviderConfig() + testMspTenantDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_msp_managed_tenant.test", "name", testMspTenantDataSource.Name),
					resource.TestCheckResourceAttr("data.cdo_msp_managed_tenant.test", "display_name", testMspTenantDataSource.DisplayName),
					resource.TestCheckResourceAttr("data.cdo_msp_managed_tenant.test", "id", testMspTenantDataSource.Id),
					resource.TestCheckResourceAttr("data.cdo_msp_managed_tenant.test", "region", testMspTenantDataSource.Region),
				),
			},
		},
	})
}
