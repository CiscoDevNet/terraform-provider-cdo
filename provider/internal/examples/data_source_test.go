package examples_test

import (
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var dataSourceModel = struct {
	Name string
}{
	Name: "example name", // typically we get the actual value from the environment like this: acctest.Env.ExampleDataSourceName(), so that most parameters are configurable so that it can be run in different CDO environment
}
var dataSourceTemplate = `
data "cdo_example" "test" {
	name = "{{.Name}}"
}`
var dataSourceConfig = acctest.MustParseTemplate(dataSourceTemplate, dataSourceModel)

func TestAccExampleDataSource(t *testing.T) {
	t.Skip("This is an example data source test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 acctest.PreCheckFunc(t),
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: acctest.ProviderConfig() + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cdo_example.test", "name", dataSourceModel.Name),
				),
			},
		},
	})
}
