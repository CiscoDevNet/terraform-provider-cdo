package ftd_test

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"
	"strconv"
	"strings"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type ResourceType struct {
	Name             string
	AccessPolicyName string
	PerformanceTier  string
	Virtual          string
	Licenses         string
	AccessPolicyUid  string
	Labels           string
}

const ResourceTemplate = `
resource "cdo_ftd_device" "test" {
	name = "{{.Name}}"
	access_policy_name = "{{.AccessPolicyName}}"
	performance_tier = "{{.PerformanceTier}}"
	virtual = "{{.Virtual}}"
	licenses = {{.Licenses}}
	labels = {{.Labels}}
}`

var testResource = ResourceType{
	Name:             acctest.Env.FtdResourceName(),
	AccessPolicyName: acctest.Env.FtdResourceAccessPolicyName(),
	PerformanceTier:  acctest.Env.FtdResourcePerformanceTier(),
	Virtual:          acctest.Env.FtdResourceVirtual(),
	Licenses:         acctest.Env.FtdResourceLicenses(),
	Labels:           acctest.Env.FtdResourceTags().GetLabelsJsonArrayString(),
}
var testResourceConfig = acctest.MustParseTemplate(ResourceTemplate, testResource)

var testResource_NewName = acctest.MustOverrideFields(testResource, map[string]any{
	"Name": acctest.Env.FtdResourceNewName(),
})

var testResourceConfig_NewName = acctest.MustParseTemplate(ResourceTemplate, testResource_NewName)

var reorderedLabels = tags.New(sliceutil.Reverse[string](tags.MustParseJsonArrayString(testResource.Labels))...).GetLabelsJsonArrayString()

var testResource_ReorderLabels = acctest.MustOverrideFields(testResource, map[string]any{
	"Labels": reorderedLabels,
})

var testResourceConfig_ReorderLabels = acctest.MustParseTemplate(ResourceTemplate, testResource_ReorderLabels)

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
					resource.TestCheckResourceAttrSet("cdo_ftd_device.test", "licenses.0"),   // there is something at position 0 of licenses array
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "licenses.#", "1"), // number of licenses = 1
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "access_policy_name", testResource.AccessPolicyName),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "labels.#", strconv.Itoa(len(acctest.Env.FtdResourceTags().Labels))),
					resource.TestCheckResourceAttrWith("cdo_ftd_device.test", "labels.0", testutil.CheckEqual(acctest.Env.FtdResourceTags().Labels[0])),
					resource.TestCheckResourceAttrWith("cdo_ftd_device.test", "labels.1", testutil.CheckEqual(acctest.Env.FtdResourceTags().Labels[1])),
					resource.TestCheckResourceAttrWith("cdo_ftd_device.test", "labels.2", testutil.CheckEqual(acctest.Env.FtdResourceTags().Labels[2])),
					resource.TestCheckResourceAttrWith("cdo_ftd_device.test", "generated_command", func(value string) error {
						ok := strings.HasPrefix(value, "configure manager add")
						if !ok {
							return fmt.Errorf("generated command should starts with \"configure manager add\", but it was \"%s\"", value)
						}
						return nil
					}),
				),
			},
			// Update order of label testing
			{
				Config:   acctest.ProviderConfig() + testResourceConfig_ReorderLabels,
				PlanOnly: true, // this will check the plan is empty
			},
			// Update Name testing
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
