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
var groupedLabels = map[string][]string{"acceptancetest": labels}

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
}`

var testIosResource = testIosResourceType{
	Name:              acctest.Env.IosResourceName(),
	SocketAddress:     acctest.Env.IosResourceSocketAddress(),
	Username:          acctest.Env.IosResourceUsername(),
	Password:          acctest.Env.IosResourcePassword(),
	ConnectorName:     acctest.Env.IosResourceConnectorName(),
	IgnoreCertificate: acctest.Env.IosResourceIgnoreCertificate(),
	Labels:            testutil.MustJson(labels),
	GroupedLabels:     testutil.MustJson(groupedLabels),

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

var testIosResource_ReplaceGroupTags = acctest.MustOverrideFields(testIosResource, map[string]any{
	"GroupedLabels": map[string][]string{
		"my-cool-new-label-group": labels,
	},
})
var testIosResourceConfig_ReplaceGroupTags = acctest.MustParseTemplate(testIosDataSourceTemplate, testIosResource_ReplaceGroupTags)

func TestAccIosDeviceResource_SDC(t *testing.T) {
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
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "labels.0", testutil.CheckEqual(labels[0])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "labels.1", testutil.CheckEqual(labels[1])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "labels.2", testutil.CheckEqual(labels[2])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "grouped_labels.acceptancetest.0", testutil.CheckEqual(labels[0])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "grouped_labels.acceptancetest.1", testutil.CheckEqual(labels[1])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "grouped_labels.acceptancetest.2", testutil.CheckEqual(labels[2])),
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
					resource.TestCheckResourceAttr("cdo_ios_device.test", "name", testIosResource_NewName.Name),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "grouped_labels.my-cool-new-label-group.0", testutil.CheckEqual(labels[0])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "grouped_labels.my-cool-new-label-group.1", testutil.CheckEqual(labels[1])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "grouped_labels.my-cool-new-label-group.2", testutil.CheckEqual(labels[2])),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
