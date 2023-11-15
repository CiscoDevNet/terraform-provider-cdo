package seconboarding_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var resourceModel = struct{}{}

const resourceTemplate = `
resource "cdo_sec" "test" {
}`

var resourceConfig = acctest.MustParseTemplate(resourceTemplate, resourceModel)

func TestAccSecResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: acctest.ProviderConfig() + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("cdo_sec.test", "name", func(value string) error {
						if !strings.ContainsAny(value, "SEC") {
							return fmt.Errorf("SEC name does not contain \"SEC\" after apply, this is likely an error")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("cdo_sec.test", "cdo_bootstrap_data", func(value string) error {
						if value == "" {
							return fmt.Errorf("CDO bootstrap data is empty after apply")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("cdo_sec.test", "sec_bootstrap_data", func(value string) error {
						if value == "" {
							return fmt.Errorf("SEC bootstrap data is empty after apply")
						}
						return nil
					}),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
