// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connector

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *cdoClient.Client
}

type DataSourceModel struct {
	Uid       types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	TenantUid types.String `tfsdk:"tenant_uid"`
	PublicKey *PublicKey   `tfsdk:"public_key"`
}
type PublicKey struct {
	EncodedKey types.String `tfsdk:"encoded_key"`
	Version    types.Int64  `tfsdk:"version"`
	KeyId      types.String `tfsdk:"key_id"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SDC data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the Secure Device Connector (SDC).",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Secure Device Connector (SDC).",
				Required:            true,
			},
			"tenant_uid": schema.StringAttribute{
				MarkdownDescription: "The tenant uid this SDC belongs to.",
				Computed:            true,
			},
			"public_key": schema.ObjectAttribute{
				MarkdownDescription: "public key of the Secure Device Connector (SDC).",
				Computed:            true,
				AttributeTypes: map[string]attr.Type{
					"encoded_key": types.StringType,
					"version":     types.Int64Type,
					"key_id":      types.StringType,
				},
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

	res, err := d.client.ReadConnectorByName(ctx, *connector.NewReadByNameInput(planData.Name.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read sdc devices", err.Error())
		return
	}

	planData.Uid = types.StringValue(res.Uid)
	planData.Name = types.StringValue(res.Name)
	planData.TenantUid = types.StringValue(res.TenantUid)

	planData.PublicKey = &PublicKey{
		EncodedKey: types.StringValue(res.PublicKey.EncodedKey),
		Version:    types.Int64Value(res.PublicKey.Version),
		KeyId:      types.StringValue(res.PublicKey.KeyId),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
