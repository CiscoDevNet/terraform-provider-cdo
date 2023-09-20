package tenant_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testTenant = struct {
	Name              string
	HumanReadableName string
	SubscriptionType  string
}{
	Name:              acctest.Env.TenantDataSourceName(),
	HumanReadableName: acctest.Env.TenantDataSourceHumanReadableName(),
	SubscriptionType:  acctest.Env.TenantDataSourceSubscriptionType(),
}

const testTenantTemplate = `
data "cdo_tenant" "test" {}`

var testTenantConfig = acctest.MustParseTemplate(testTenantTemplate, testTenant)

func TestAccTenantDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + testTenantConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_tenant.test", "name", testTenant.Name),
					resource.TestCheckResourceAttr("data.cdo_tenant.test", "human_readable_name", testTenant.HumanReadableName),
					resource.TestCheckResourceAttr("data.cdo_tenant.test", "subscription_type", testTenant.SubscriptionType),
				),
			},
		},
	})
}
