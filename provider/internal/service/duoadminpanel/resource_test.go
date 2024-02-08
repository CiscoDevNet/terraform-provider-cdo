package duoadminpanel_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/testutil"
	"strconv"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var resourceModel = struct {
	Name           string
	Host           string
	IntegrationKey string
	SecretKey      string
	Labels         string
}{
	Name:           acctest.Env.DuoAdminPanelResourceName(),
	Host:           acctest.Env.DuoAdminPanelResourceHost(),
	IntegrationKey: acctest.Env.DuoAdminPanelResourceIntegrationKey(),
	SecretKey:      acctest.Env.DuoAdminPanelResourceSecretKey(),
	Labels:         acctest.Env.DuoAdminPanelResourceTags().GetLabelsJsonArrayString(),
}

const resourceTemplate = `
resource "cdo_duo_admin_panel" "test" {
	name = "{{.Name}}"
	host = "{{.Host}}"
	integration_key = "{{.IntegrationKey}}"
	secret_key = "{{.SecretKey}}"
	labels = {{.Labels}}
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
					resource.TestCheckResourceAttr("cdo_duo_admin_panel.test", "labels.#", strconv.Itoa(len(acctest.Env.DuoAdminPanelResourceTags().Labels))),
					resource.TestCheckResourceAttrWith("cdo_duo_admin_panel.test", "labels.0", testutil.CheckEqual(acctest.Env.DuoAdminPanelResourceTags().Labels[0])),
					resource.TestCheckResourceAttrWith("cdo_duo_admin_panel.test", "labels.1", testutil.CheckEqual(acctest.Env.DuoAdminPanelResourceTags().Labels[1])),
					resource.TestCheckResourceAttrWith("cdo_duo_admin_panel.test", "labels.2", testutil.CheckEqual(acctest.Env.DuoAdminPanelResourceTags().Labels[2])),
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
