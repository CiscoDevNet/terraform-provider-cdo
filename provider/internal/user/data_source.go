package user

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceModel struct {
	Uid         types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ApiOnlyUser types.Bool   `tfsdk:"is_api_only_user"`
	UserRole    types.String `tfsdk:"role"`
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *cdoClient.Client
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get the identifiers of users to be referenced elsewhere, e.g., to add an existing user to a new tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier for the user.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the user.",
				Required:            true,
			},
			"is_api_only_user": schema.BoolAttribute{
				MarkdownDescription: "CDO has two kinds of users: actual users with email addresses and API-only users for programmatic access. This boolean indicates what type of user this is.",
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Roles assigned to the user in this tenant.",
				Computed:            true,
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var planData DataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.ReadUserByUsername(ctx, *user.NewReadByUsernameInput(planData.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user", err.Error())
		return
	}

	planData.Uid = types.StringValue(res.Uid)
	planData.Name = types.StringValue(res.Name)
	planData.ApiOnlyUser = types.BoolValue(res.ApiOnlyUser)
	planData.UserRole = types.StringValue(res.UserRoles[0]) // while in theory the API supports multiple roles, our UI restricts users to one role. So we're doing the same here
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
