package ftd_test

import (
	"fmt"
	"strings"
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
	Name:             "test-cloud-ftd-12345",
	AccessPolicyName: "Default Access Control Policy",
	PerformanceTier:  "FTDv5",
	Virtual:          "true",
	Licenses:         "[\"BASE\"]",
}
var testResourceConfig = acctest.MustParseTemplate(ResourceTemplate, testResource)

var testResource_NewName = acctest.MustOverrideFields(testResource, map[string]any{
	"Name": "test-cloud-ftd-new-name",
})

var testResourceConfig_NewName = acctest.MustParseTemplate(ResourceTemplate, testResource_NewName)

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
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "licenses.#", "1"), // number of licenses = 1
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "access_policy_name", testResource.AccessPolicyName),
					resource.TestCheckResourceAttrWith("cdo_ftd_device.test", "generated_command", func(value string) error {
						ok := strings.HasPrefix(value, "configure manager add")
						if !ok {
							return fmt.Errorf("generated command should starts with \"configure manager add\", but it was \"%s\"", value)
						}
						return nil
					}),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testResourceConfig_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "name", testResource_NewName.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
