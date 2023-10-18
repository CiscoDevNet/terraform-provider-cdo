package ios_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"
	"strconv"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type testIosResourceType struct {
	Name              string
	SocketAddress     string
	ConnectorType     string
	Username          string
	Password          string
	ConnectorName     string
	IgnoreCertificate string
	Labels            string

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
	Labels:            acctest.Env.IosResourceTags().GetLabelsJsonArrayString(),

	Host: acctest.Env.IosResourceHost(),
	Port: acctest.Env.IosResourcePort(),
}
var testIosResourceConfig = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource)

var reorderedLabels = tags.New(sliceutil.Reverse[string](tags.MustParseJsonArrayString(testIosResource.Labels))...).GetLabelsJsonArrayString()

var testIosResource_ReorderedLabels = acctest.MustOverrideFields(testIosResource, map[string]any{
	"Labels": reorderedLabels,
})
var testIosResourceConfig_ReorderedLabels = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource_ReorderedLabels)

var testIosResource_NewName = acctest.MustOverrideFields(testIosResource, map[string]any{
	"Name": acctest.Env.IosResourceNewName(),
})
var testIosResourceConfig_NewName = acctest.MustParseTemplate(testIosResourceTemplate, testIosResource_NewName)

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
					resource.TestCheckResourceAttr("cdo_ios_device.test", "labels.#", strconv.Itoa(len(acctest.Env.FtdResourceTags().Labels))),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "labels.0", testutil.CheckEqual(acctest.Env.IosResourceTags().Labels[0])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "labels.1", testutil.CheckEqual(acctest.Env.IosResourceTags().Labels[1])),
					resource.TestCheckResourceAttrWith("cdo_ios_device.test", "labels.2", testutil.CheckEqual(acctest.Env.IosResourceTags().Labels[2])),
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
			// Delete testing automatically occurs in TestCase
		},
	})
}
