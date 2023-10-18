// Package asa_test does not contain CDG test as we do not want to be using ASAs accessible from the public subnet for our tests
package asa_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"regexp"
	"strconv"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type testAsaResourceType struct {
	Name              string
	SocketAddress     string
	ConnectorName     string
	ConnectorType     string
	Username          string
	Password          string
	IgnoreCertificate bool
	Labels            string

	Host string
	Port int64
}

const asaResourceTemplate = `
resource "cdo_asa_device" "test" {
	name = "{{.Name}}"
	socket_address = "{{.SocketAddress}}"
	connector_name = "{{.ConnectorName}}"
	connector_type = "{{.ConnectorType}}"
	username = "{{.Username}}"
	password = "{{.Password}}"
	ignore_certificate = "{{.IgnoreCertificate}}"
	labels = {{.Labels}}
}`

const asaResourceTemplateNoLabels = `
resource "cdo_asa_device" "test" {
	name = "{{.Name}}"
	socket_address = "{{.SocketAddress}}"
	connector_name = "{{.ConnectorName}}"
	connector_type = "{{.ConnectorType}}"
	username = "{{.Username}}"
	password = "{{.Password}}"
	ignore_certificate = "{{.IgnoreCertificate}}"
}`

// SDC configs.

// default config.
var testAsaResource_SDC = testAsaResourceType{
	Name:              acctest.Env.AsaResourceSdcName(),
	SocketAddress:     acctest.Env.AsaResourceSdcSocketAddress(),
	ConnectorName:     acctest.Env.AsaResourceSdcConnectorName(),
	ConnectorType:     acctest.Env.AsaResourceSdcConnectorType(),
	Username:          acctest.Env.AsaResourceSdcUsername(),
	Password:          acctest.Env.AsaResourceSdcPassword(),
	IgnoreCertificate: acctest.Env.AsaResourceSdcIgnoreCertificate(),
	Labels:            acctest.Env.AsaResourceSdcTags().GetLabelsJsonArrayString(),

	Host: acctest.Env.AsaResourceSdcHost(),
	Port: acctest.Env.AsaResourceSdcPort(),
}

var testAsaResourceConfig_SDC = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC)
var testAsaResourceConfig_SDC_NoLabels = acctest.MustParseTemplate(asaResourceTemplateNoLabels, testAsaResource_SDC)

// new label order config.
var reorderedLabels = tags.New(sliceutil.Reverse[string](tags.MustParseJsonArrayString(testAsaResource_SDC.Labels))...).GetLabelsJsonArrayString()

var testAsaResource_SDC_ReorderedLabels = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{
	"Labels": reorderedLabels,
})
var testAsaResourceConfig_SDC_ReorderedLabels = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_ReorderedLabels)

// new name config.
var testAsaResource_SDC_NewName = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{
	"Name": acctest.Env.AsaResourceSdcNewName(),
})
var testAsaResourceConfig_SDC_NewName = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_NewName)

// new creds config.
var testAsaResource_SDC_BadCreds = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{
	"Password": acctest.Env.AsaResourceSdcWrongPassword(),
})
var testAsaResourceConfig_SDC_NewCreds = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_BadCreds)

func TestAccAsaDeviceResource_SDC(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_SDC.SocketAddress),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "host", testAsaResource_SDC.Host),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "port", strconv.FormatInt(testAsaResource_SDC.Port, 10)),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "connector_type", testAsaResource_SDC.ConnectorType),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "username", testAsaResource_SDC.Username),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_SDC.Password),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "labels.#", strconv.Itoa(len(acctest.Env.AsaResourceSdcTags().Labels))),
					resource.TestCheckResourceAttrWith("cdo_asa_device.test", "labels.0", testutil.CheckEqual(acctest.Env.AsaResourceSdcTags().Labels[0])),
					resource.TestCheckResourceAttrWith("cdo_asa_device.test", "labels.1", testutil.CheckEqual(acctest.Env.AsaResourceSdcTags().Labels[1])),
					resource.TestCheckResourceAttrWith("cdo_asa_device.test", "labels.2", testutil.CheckEqual(acctest.Env.AsaResourceSdcTags().Labels[2])),
				),
			},
			// bad credential tests
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC_BadCreds.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_SDC_BadCreds.Password),
				),
				ExpectError: regexp.MustCompile(`.*bad credentials.*`),
			},
			// fix correct credentials test
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_SDC.Password),
				),
			},
			// Update order of label testing
			{
				Config:   acctest.ProviderConfig() + testAsaResourceConfig_SDC_ReorderedLabels,
				PlanOnly: true, // this will check the plan is empty
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC_NewName.Name),
				),
			},

			// change location test - disabled until we create another asa
			// {
			// 	Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewLocation,
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_SDC_NewLocation.SocketAddress),
			// 	),
			// },
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAsaDeviceResource_SDC_NoLabels(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NoLabels,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "labels.#", strconv.Itoa(0)),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
