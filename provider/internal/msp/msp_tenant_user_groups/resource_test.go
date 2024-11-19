package msp_tenant_user_groups_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

type UserGroup struct {
	GroupIdentifier string
	IssuerUrl       string
	Name            string
	Role            string
	Notes           string
}

var testMspManagedTenantUserGroupsResource = struct {
	TenantUid  string
	UserGroups []UserGroup
}{
	UserGroups: []UserGroup{
		{GroupIdentifier: "customer-managers", IssuerUrl: "https://www.customer-idp.com", Name: "developers", Role: "ROLE_SUPER_ADMIN", Notes: "Managers in customer's organization"},
		{GroupIdentifier: "msp-managers", IssuerUrl: "https://www.msp-idp.com", Name: "managers", Role: "ROLE_ADMIN", Notes: "Managers in MSP organization"},
	},
	TenantUid: acctest.Env.MspTenantId(),
}

const testMspManagedTenantUserGroupsTemplate = `
resource "cdo_msp_managed_tenant_user_groups" "test" {
	tenant_uid = "{{.TenantUid}}"
	user_groups = [
		{
			"group_identifier": "{{(index .UserGroups 0).GroupIdentifier }}"
			"issuer_url": "{{ (index .UserGroups 0).IssuerUrl }}"
			"name": "{{ (index .UserGroups 0).Name }}"
			"role": "{{ (index .UserGroups 0).Role }}"
			"notes": "{{ (index .UserGroups 0).Notes }}"
		},
		{
			"group_identifier": "{{(index .UserGroups 1).GroupIdentifier }}"
			"issuer_url": "{{ (index .UserGroups 1).IssuerUrl }}"
			"name": "{{ (index .UserGroups 1).Name }}"
			"role": "{{ (index .UserGroups 1).Role }}"
			"notes": "{{ (index .UserGroups 1).Notes }}"
		}
	]
}`

var testMspManagedTenantUserGroupsResourceConfig = acctest.MustParseTemplate(testMspManagedTenantUserGroupsTemplate, testMspManagedTenantUserGroupsResource)

func TestAccMspManagedTenantUserGroupsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.MspProviderConfig() + testMspManagedTenantUserGroupsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"tenant_uid",
						testMspManagedTenantUserGroupsResource.TenantUid),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"user_groups.0.group_identifier",
						testMspManagedTenantUserGroupsResource.UserGroups[0].GroupIdentifier),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"user_groups.0.issuer_url",
						testMspManagedTenantUserGroupsResource.UserGroups[0].IssuerUrl),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"user_groups.0.name",
						testMspManagedTenantUserGroupsResource.UserGroups[0].Name),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"user_groups.0.notes",
						testMspManagedTenantUserGroupsResource.UserGroups[0].Notes),
					resource.TestCheckResourceAttr(
						"cdo_msp_managed_tenant_user_groups.test",
						"user_groups.0.role",
						testMspManagedTenantUserGroupsResource.UserGroups[0].Role,
					),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"user_groups.1.group_identifier",
						testMspManagedTenantUserGroupsResource.UserGroups[1].GroupIdentifier),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"user_groups.1.issuer_url",
						testMspManagedTenantUserGroupsResource.UserGroups[1].IssuerUrl),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"user_groups.1.name",
						testMspManagedTenantUserGroupsResource.UserGroups[1].Name),
					resource.TestCheckResourceAttr("cdo_msp_managed_tenant_user_groups.test",
						"user_groups.1.notes",
						testMspManagedTenantUserGroupsResource.UserGroups[1].Notes),
					resource.TestCheckResourceAttr(
						"cdo_msp_managed_tenant_user_groups.test",
						"user_groups.1.role",
						testMspManagedTenantUserGroupsResource.UserGroups[1].Role,
					),
				),
			},
		},
	})
}
