// Package examples package provides an example for creating terraform resource and data source
package examples

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = NewExampleDataSource()

// NewExampleDataSource is used in provider.go to register this data source with terraform.
func NewExampleDataSource() datasource.DataSource {
	return &ExampleDataSource{}
}

// ExampleDataSource is the struct object that will be consumed by terraform, it contains methods that defines metadata, schema, read, create, etc...
type ExampleDataSource struct {
	client *cdoClient.Client
}

/////
// model class: define mapping from Go types to Terraform types.
/////

type ExampleDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Metadata is primarily used to define the name for this data source.
func (d *ExampleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ios_device"
}

// Schema is primarily used to define the terraform schema for this data source.
// it needs to be consistent with the Model classes' `tfsdk:"xxx"` above.
func (d *ExampleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example Id description",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Example name description",
				Required:            true,
			},
		},
	}
}

// Configure sets the environment for CRUD operations, i.e. we set the CDO Go client to use when doing the CRUD operations.
// The Go client is configured in the provider, and it is available here, so we just set it directly.
func (d *ExampleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdoClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cdoClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read operation is all it needs for terraform data source.
// this function is responsible for:
// 1. mapping the cdo go client's data to the model classes we defined above.
// 2. report any error using `resp.Diagnostics`.
// 3. set model classes to the terraform state.
func (r *ExampleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Trace(ctx, "reading Example data source")

	// 1. read state data
	var stateData ExampleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	if err := ReadDataSource(ctx, r, &stateData); err != nil {
		resp.Diagnostics.AddError("failed to read Example data source", err.Error())
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}
