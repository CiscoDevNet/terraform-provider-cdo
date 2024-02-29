package duoadminpanel_test

import (
	"strconv"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var labels = []string{"testdevice", "acceptancetest", "terraformprovider"}
var groupedLabels = map[string][]string{"acceptancetest": sliceutil.Map(labels, func(input string) string { return "grouped-" + input })}

var resourceModel = struct {
	Name           string
	Host           string
	IntegrationKey string
	SecretKey      string
	Labels         string
	GroupedLabels  string
}{
	Name:           acctest.Env.DuoAdminPanelResourceName(),
	Host:           acctest.Env.DuoAdminPanelResourceHost(),
	IntegrationKey: acctest.Env.DuoAdminPanelResourceIntegrationKey(),
	SecretKey:      acctest.Env.DuoAdminPanelResourceSecretKey(),
	Labels:         testutil.MustJson(labels),
	GroupedLabels:  acctest.MustGenerateLabelsTF(groupedLabels),
}

const resourceTemplate = `
resource "cdo_duo_admin_panel" "test" {
	name = "{{.Name}}"
	host = "{{.Host}}"
	integration_key = "{{.IntegrationKey}}"
	secret_key = "{{.SecretKey}}"
	labels = {{.Labels}}
	grouped_labels = {{.GroupedLabels}}
}`

var resourceConfig = acctest.MustParseTemplate(resourceTemplate, resourceModel)

var resourceModel_NewName = acctest.MustOverrideFields(resourceModel, map[string]any{
	"Name": acctest.Env.DuoAdminPanelResourceNewName(),
})
var resourceConfig_NewName = acctest.MustParseTemplate(resourceTemplate, resourceModel_NewName)

func TestAccDuoAdminPanelResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "name", resourceModel.Name),
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "host", resourceModel.Host),
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "integration_key", resourceModel.IntegrationKey),
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "secret_key", resourceModel.SecretKey),
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "labels.#", strconv.Itoa(len(labels))),
					resource.TestCheckTypeSetElemAttr("cdo_duo_admin_panel.test", "labels.*", labels[0]),
					resource.TestCheckTypeSetElemAttr("cdo_duo_admin_panel.test", "labels.*", labels[1]),
					resource.TestCheckTypeSetElemAttr("cdo_duo_admin_panel.test", "labels.*", labels[2]),
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "grouped_labels.%", "1"),
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "grouped_labels.acceptancetest.#", strconv.Itoa(len(groupedLabels["acceptancetest"]))),
					resource.TestCheckTypeSetElemAttr("cdo_duo_admin_panel.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][0]),
					resource.TestCheckTypeSetElemAttr("cdo_duo_admin_panel.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][1]),
					resource.TestCheckTypeSetElemAttr("cdo_duo_admin_panel.test", "grouped_labels.acceptancetest.*", groupedLabels["acceptancetest"][2]),
				),
			},
			// Update and Read testing
			{
				Config: acctest.ProviderConfig() + resourceConfig_NewName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "name", resourceModel_NewName.Name),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
