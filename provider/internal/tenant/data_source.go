package tenant

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type DataSourceModel struct {
	Uid               types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	HumanReadableName types.String `tfsdk:"human_readable_name"`
	SubscriptionType  types.String `tfsdk:"subscription_type"`
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *cdoClient.Client
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information on the tenant upon which the Terraform provider is performing operations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier for the tenant.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the tenant.",
				Computed:            true,
			},
			"human_readable_name": schema.StringAttribute{
				MarkdownDescription: "Human-readable name of the tenant as displayed on the CDO UI (if different from the tenant name).",
				Computed:            true,
			},
			"subscription_type": schema.StringAttribute{
				MarkdownDescription: "The type of CDO subscription used on this tenant.",
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

	res, err := d.client.ReadTenantDetails(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read tenant", err.Error())
		return
	}

	planData.Uid = types.StringValue(res.UserAuthentication.Details.TenantUid)
	planData.Name = types.StringValue(res.UserAuthentication.Details.TenantName)
	planData.HumanReadableName = types.StringValue(res.UserAuthentication.Details.TenantOrganizationName)
	planData.SubscriptionType = types.StringValue(res.UserAuthentication.Details.TenantPayType)
	tflog.Debug(ctx, fmt.Sprintf("Read tenant details %+v", planData))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
