package cdfmc

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceModel struct {
	Uid             types.String `tfsdk:"id"`
	Hostname        types.String `tfsdk:"hostname"`
	SoftwareVersion types.String `tfsdk:"software_version"`
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *cdoClient.Client
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cdfmc"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information on the cloud-delivered FMC in your tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier for the cdFMC.",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Name of the tenant.",
				Computed:            true,
			},
			"software_version": schema.StringAttribute{
				MarkdownDescription: "Software version of the cdFMC.",
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

	cloudFmc, err := d.client.ReadCloudFmc(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read cdFMC", err.Error())
		return
	}
	planData.Uid = types.StringValue(cloudFmc.Uid)
	planData.Hostname = types.StringValue(cloudFmc.Host)
	planData.SoftwareVersion = types.StringValue(cloudFmc.SoftwareVersion)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
