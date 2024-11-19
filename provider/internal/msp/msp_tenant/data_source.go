package msp_tenant

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/tenants"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewTenantDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *cdoClient.Client
}

func (d *DataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_msp_managed_tenant"
}

func (d *DataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information on an MSP-managed tenant in your portal.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the tenant",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the tenant",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the tenant",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "CDO region in which the tenant is created. This is the same region as the region of the MSP portal.",
				Computed:            true,
			},
		},
	}
}

func (d *DataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var planData TenantDatasourceModel

	// Read Terraform configuration data into the model
	response.Diagnostics.Append(request.Config.Get(ctx, &planData)...)
	if response.Diagnostics.HasError() {
		return
	}

	mspManagedTenants, err := d.client.FindMspManagedTenantByName(ctx, tenants.ReadByNameInput{
		Name: planData.Name.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError("Failed to read MSP Managed Tenant", err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Found %d MSP managed tenants by name %s", mspManagedTenants.Count, planData.Name.ValueString()))
	if mspManagedTenants.Count != 1 {
		response.Diagnostics.AddError("Cannot find MSP managed tenant by name "+planData.Name.ValueString(), fmt.Sprintf("Found %d tenants by name %s", mspManagedTenants.Count, planData.Name.ValueString()))
		return
	}

	mspManagedTenant := mspManagedTenants.Items[0]
	planData.Id = types.StringValue(mspManagedTenant.Uid)
	planData.Name = types.StringValue(mspManagedTenant.Name)
	planData.DisplayName = types.StringValue(mspManagedTenant.DisplayName)
	planData.Region = types.StringValue(mspManagedTenant.Region)

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &planData)...)
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
