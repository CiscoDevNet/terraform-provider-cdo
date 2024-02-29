package ftd_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var labels = []string{"acceptancetest", "terraform", "test_ftd"}
var groupedLabels = map[string][]string{"acceptancetest": sliceutil.Map(labels, func(input string) string { return "grouped-" + input })}

type ResourceType struct {
	Name             string
	AccessPolicyName string
	PerformanceTier  string
	Virtual          string
	Licenses         string
	AccessPolicyUid  string
	Labels           string
	GroupedLabels    string
}

const ResourceTemplate = `
resource "cdo_ftd_device" "test" {
	name = "{{.Name}}"
	access_policy_name = "{{.AccessPolicyName}}"
	performance_tier = "{{.PerformanceTier}}"
	virtual = "{{.Virtual}}"
	licenses = {{.Licenses}}
	labels = {{.Labels}}
	grouped_labels = {{.GroupedLabels}}
}`

var testResource = ResourceType{
	Name:             acctest.Env.FtdResourceName(),
	AccessPolicyName: acctest.Env.FtdResourceAccessPolicyName(),
	PerformanceTier:  acctest.Env.FtdResourcePerformanceTier(),
	Virtual:          acctest.Env.FtdResourceVirtual(),
	Licenses:         acctest.Env.FtdResourceLicenses(),
	Labels:           testutil.MustJson(labels),
	GroupedLabels:    acctest.MustGenerateLabelsTF(groupedLabels),
}
var testResourceConfig = acctest.MustParseTemplate(ResourceTemplate, testResource)

var testResource_NewName = acctest.MustOverrideFields(testResource, map[string]any{
	"Name": acctest.Env.FtdResourceNewName(),
})

var testResourceConfig_NewName = acctest.MustParseTemplate(ResourceTemplate, testResource_NewName)

var reorderedLabels = testutil.MustJson(sliceutil.Reverse(labels))

var testResource_ReorderLabels = acctest.MustOverrideFields(testResource, map[string]any{
	"Labels": reorderedLabels,
})

var testResourceConfig_ReorderLabels = acctest.MustParseTemplate(ResourceTemplate, testResource_ReorderLabels)

var renamedGroupedLabels = map[string][]string{
	"my-cool-new-label-group": groupedLabels["acceptancetest"],
}
var testResource_ReplaceGroupedLabels = acctest.MustOverrideFields(testResource, map[string]any{
	"GroupedLabels": acctest.MustGenerateLabelsTF(renamedGroupedLabels),
})
var testResourceConfig_ReplaceGroupedLabels = acctest.MustParseTemplate(ResourceTemplate, testResource_ReplaceGroupedLabels)

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
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "labels.#", strconv.Itoa(len(labels))),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "labels.*", labels[0]),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "labels.*", labels[1]),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "labels.*", labels[2]),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "grouped_labels.%", "1"),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "grouped_labels.acceptancetest.#", strconv.Itoa(len(groupedLabels["acceptancetest"]))),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][0]),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][1]),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][2]),
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
			// Replace Grouped Labels
			{
				Config: acctest.ProviderConfig() + testResourceConfig_ReplaceGroupedLabels,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "grouped_labels.%", "1"),
					resource.TestCheckResourceAttr("cdo_ftd_device.test", "grouped_labels.my-cool-new-label-group.#", strconv.Itoa(len(renamedGroupedLabels["my-cool-new-label-group"]))),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "grouped_labels.my-cool-new-label-group.*", renamedGroupedLabels["my-cool-new-label-group"][0]),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "grouped_labels.my-cool-new-label-group.*", renamedGroupedLabels["my-cool-new-label-group"][1]),
					resource.TestCheckTypeSetElemAttr("cdo_ftd_device.test", "grouped_labels.my-cool-new-label-group.*", renamedGroupedLabels["my-cool-new-label-group"][2]),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
