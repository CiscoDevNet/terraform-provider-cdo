package asa_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"
	"regexp"
	"strconv"
	"testing"

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

var alternativeDeviceLocation = acctest.Env.AsaResourceAlternativeDeviceLocation()

var testAsaResourceConfig_SDC = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC)

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

var testAsaResource_SDC_NewLocation = acctest.MustOverrideFields(testAsaResource_SDC, map[string]any{"SocketAddress": alternativeDeviceLocation})
var testAsaResourceConfig_SDC_NewLocation = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_SDC_NewLocation)

// CDG configs

// default config.
var testAsaResource_CDG = testAsaResourceType{
	Name:              acctest.Env.AsaResourceCdgName(),
	SocketAddress:     acctest.Env.AsaResourceCdgSocketAddress(),
	ConnectorName:     acctest.Env.AsaResourceCdgConnectorName(),
	ConnectorType:     acctest.Env.AsaResourceCdgConnectorType(),
	Username:          acctest.Env.AsaResourceCdgUsername(),
	Password:          acctest.Env.AsaResourceCdgPassword(),
	IgnoreCertificate: acctest.Env.AsaResourceCdgIgnoreCertificate(),
	Labels:            acctest.Env.AsaResourceCdgTags().GetLabelsJsonArrayString(),

	Host: acctest.Env.AsaResourceCdgHost(),
	Port: acctest.Env.AsaResourceCdgPort(),
}
var testAsaResourceConfig_CDG = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG)

// new name config.
var testAsaResource_CDG_NewName = acctest.MustOverrideFields(testAsaResource_CDG, map[string]any{
	"Name": acctest.Env.AsaResourceCdgNewName(),
})
var testAsaResourceConfig_CDG_NewName = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG_NewName)

// bad credentials config.
var testAsaResource_CDG_BadCreds = acctest.MustOverrideFields(testAsaResource_CDG, map[string]any{
	"Password": acctest.Env.AsaResourceCdgWrongPassword(),
})
var testAsaResourceConfig_CDG_NewCreds = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG_BadCreds)

var testAsaResource_CDG_NewLocation = acctest.MustOverrideFields(testAsaResource_CDG, map[string]any{"SocketAddress": alternativeDeviceLocation})
var testAsaResourceConfig_CDG_NewLocation = acctest.MustParseTemplate(asaResourceTemplate, testAsaResource_CDG_NewLocation)

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
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_SDC_NewName.Name),
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
			// change location test
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_SDC_NewLocation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_SDC_NewLocation.SocketAddress),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAsaDeviceResource_CDG(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_CDG.SocketAddress),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "host", testAsaResource_CDG.Host),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "port", strconv.FormatInt(testAsaResource_CDG.Port, 10)),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "connector_type", testAsaResource_CDG.ConnectorType),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "username", testAsaResource_CDG.Username),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_CDG.Password),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "labels.#", strconv.Itoa(len(acctest.Env.AsaResourceCdgTags().Labels))),
					resource.TestCheckResourceAttrWith("cdo_asa_device.test", "labels.0", testutil.CheckEqual(acctest.Env.AsaResourceCdgTags().Labels[0])),
					resource.TestCheckResourceAttrWith("cdo_asa_device.test", "labels.1", testutil.CheckEqual(acctest.Env.AsaResourceCdgTags().Labels[1])),
					resource.TestCheckResourceAttrWith("cdo_asa_device.test", "labels.2", testutil.CheckEqual(acctest.Env.AsaResourceCdgTags().Labels[2])),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG_NewName.Name),
				),
			},
			// bad credentials tests
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG_NewCreds,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG_BadCreds.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_CDG_BadCreds.Password),
				),
				ExpectError: regexp.MustCompile(".*bad credentials.*"),
			},
			// fix bad credentials test
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "name", testAsaResource_CDG.Name),
					resource.TestCheckResourceAttr("cdo_asa_device.test", "password", testAsaResource_CDG.Password),
				),
			},
			// change location test
			{
				Config: acctest.ProviderConfig() + testAsaResourceConfig_CDG_NewLocation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_asa_device.test", "socket_address", testAsaResource_CDG_NewLocation.SocketAddress),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
