package ftd_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type ResourceType struct {
	ID               string
	Name             string
	AccessPolicyName string
	PerformanceTier  string
	Virtual          string
	Licenses         string
	AccessPolicyUid  string
	GeneratedCommand string
}

const ResourceTemplate = `
resource "cdo_ftd_device" "test" {
	name = "{{.Name}}"
	access_policy_name = "{{.AccessPolicyName}}"
	performance_tier = "{{.PerformanceTier}}"
	virtual = "{{.Virtual}}"
	licenses = {{.Licenses}}
}`

var testResource = ResourceType{
	Name:             "ci-test-ftd-9",
	AccessPolicyName: "Default Access Control Policy",
	PerformanceTier:  "FTDv5",
	Virtual:          "false",
	Licenses:         "[\"BASE\"]",
	GeneratedCommand: "configure manager add terraform-provider-cdo.app.staging.cdo.cisco.com LvWGkKjYNrqZlYbz2JGZqbD0ibDuxlSp h2zTtFTvwxgDIbI9pGshHNWrJGDT0jzC terraform-provider-cdo.app.staging.cdo.cisco.com",
}
var testResourceConfig = acctest.MustParseTemplate(ResourceTemplate, testResource)

func TestAccFtdResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "name", testResource.Name),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "access_policy_name", testResource.AccessPolicyName),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "performance_tier", testResource.PerformanceTier),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "virtual", testResource.Virtual),
					resource.TestCheckResourceAttrSet("cdo_ftd_device.test", "licenses.0"),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "access_policy_name", testResource.AccessPolicyName),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "generated_command", testResource.GeneratedCommand),
				),
			},
			// Update and Read testing
			//{
			//	Config: acctest.ProviderConfig() + testResourceConfig,
			//	Check: resource.ComposeAggregateTestCheckFunc(
			//		resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testResource.Name),
			//	),
			//},
			// Delete testing automatically occurs in TestCase
		},
	})
}
