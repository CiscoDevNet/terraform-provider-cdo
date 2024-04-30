package ios_test

import (
	"strconv"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var labels = []string{"acceptancetest", "test-ios-device", "terraform"}
var groupedLabels = map[string][]string{"acceptancetest": sliceutil.Map(labels, func(input string) string { return "grouped-" + input })}

type testIosResourceType struct {
	Name              string
	SocketAddress     string
	ConnectorType     string
	Username          string
	Password          string
	ConnectorName     string
	IgnoreCertificate string
	Labels            string
	GroupedLabels     string

	Host string
	Port int64
}

const testIosResourceTemplate = `
resource "cdo_ios_device" "test" {
	name = "{{.Name}}"
	socket_address = "{{.SocketAddress}}"
	username = "{{.Username}}"
	password = "{{.Password}}"
	connector_name = "{{.ConnectorName}}"
	ignore_certificate = "{{.IgnoreCertificate}}"
	labels = {{.Labels}}
	grouped_labels = {{.GroupedLabels}}
}`

var testIosResource = testIosResourceType{
	Name:              acctest.Env.IosResourceName(),
	SocketAddress:     acctest.Env.IosResourceSocketAddress(),
	Username:          acctest.Env.IosResourceUsername(),
	Password:          acctest.Env.IosResourcePassword(),
	ConnectorName:     acctest.Env.IosResourceConnectorName(),
	IgnoreCertificate: acctest.Env.IosResourceIgnoreCertificate(),
	Labels:            testutil.MustJson(labels),
	GroupedLabels:     acctest.MustGenerateLabelsTF(groupedLabels),

	Host: acctest.Env.IosResourceHost(),
	Port: acctest.Env.IosResourcePort(),
}
var testIosResourceConfig = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource)

var reorderedLabels = testutil.MustJson(sliceutil.Reverse(labels))

var testIosResource_ReorderedLabels = acctest.MustOverrideFields(testIosResource, map[string]any{
	"Labels": reorderedLabels,
})
var testIosResourceConfig_ReorderedLabels = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource_ReorderedLabels)

var testIosResource_NewName = acctest.MustOverrideFields(testIosResource, map[string]any{
	"Name": acctest.Env.IosResourceNewName(),
})
var testIosResourceConfig_NewName = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource_NewName)

var renamedGroupedLabels = map[string][]string{
	"my-cool-new-label-group": groupedLabels["acceptancetest"],
}
var testIosResource_ReplaceGroupTags = acctest.MustOverrideFields(testIosResource, map[string]any{
	"GroupedLabels": acctest.MustGenerateLabelsTF(renamedGroupedLabels),
})
var testIosResourceConfig_ReplaceGroupTags = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource_ReplaceGroupTags)

func TestAccIosDeviceResource_SDC(t *testing.T) {
	t.Skip("Disabling this test because the vSphere SDC lab is down. Yay!")
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testIosResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testIosResource.Name),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "socket_address", testIosResource.SocketAddress),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "host", testIosResource.Host),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "port", strconv.FormatInt(testIosResource.Port, 10)),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "username", testIosResource.Username),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "password", testIosResource.Password),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "labels.#", strconv.Itoa(len(labels))),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "labels.*", labels[0]),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "labels.*", labels[1]),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "labels.*", labels[2]),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "grouped_labels.%", "1"),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "grouped_labels.acceptancetest.#", strconv.Itoa(len(groupedLabels["acceptancetest"]))),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][0]),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][1]),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][2]),
				),
			},
			// Update order of label testing
			{
				Config:   acctest.ProviderConfig() + testIosResourceConfig_ReorderedLabels,
				PlanOnly: true, // this will check the plan is empty
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testIosResourceConfig_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testIosResource_NewName.Name),
				),
			},

			// Replace group labels test
			{
				Config: acctest.ProviderConfig() + testIosResourceConfig_ReplaceGroupTags,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_ios_device.test", "grouped_labels.%", "1"),
					resource.TestCheckResourceAttr("cdo_ios_device.test", "grouped_labels.my-cool-new-label-group.#", strconv.Itoa(len(renamedGroupedLabels["my-cool-new-label-group"]))),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "grouped_labels.my-cool-new-label-group.*", renamedGroupedLabels["my-cool-new-label-group"][0]),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "grouped_labels.my-cool-new-label-group.*", renamedGroupedLabels["my-cool-new-label-group"][1]),
					resource.TestCheckTypeSetElemAttr("cdo_ios_device.test", "grouped_labels.my-cool-new-label-group.*", renamedGroupedLabels["my-cool-new-label-group"][2]),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
