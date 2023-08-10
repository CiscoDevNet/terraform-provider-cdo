// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdc

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
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
	Id   types.String `tfsdk:"id"`
	Sdcs []Model      `tfsdk:"sdcs"`
}

type Model struct {
	Uid          types.String   `tfsdk:"uid"`
	Name         types.String   `tfsdk:"name"`
	TenantUid    types.String   `tfsdk:"tenant_uid"`
	SdcPublicKey PublicKeyModel `tfsdk:"sdc_public_key"`
}

type PublicKeyModel struct {
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
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"sdcs": schema.ListNestedAttribute{
				MarkdownDescription: "List of sdcs",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uid": schema.StringAttribute{
							MarkdownDescription: "Uid",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Computed:            true,
						},
						"tenant_uid": schema.StringAttribute{
							MarkdownDescription: "The tenant uid this SDC belongs to",
							Computed:            true,
						},
						"sdc_public_key": schema.ObjectAttribute{
							MarkdownDescription: "SDC public key",
							Computed:            true,
							AttributeTypes: map[string]attr.Type{
								"encoded_key": types.StringType,
								"version":     types.Int64Type,
								"key_id":      types.StringType,
							},
						},
					},
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

	var configData DataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.ReadAllSdcs(ctx, sdc.ReadAllInput{})
	if err != nil {
		resp.Diagnostics.AddError("Failed to read sdc devices", err.Error())
		return
	}

	if len(*res) == 0 {
		resp.Diagnostics.AddError("No sdc found.", "Either no sdc device is found or no default sdc is configured")
		return
	}
	configData.Id = types.StringValue((*res)[0].TenantUid)

	var sdcModels []Model
	for _, sdc := range *res {
		sdcModel := Model{
			Uid:       types.StringValue(sdc.Uid),
			Name:      types.StringValue(sdc.Name),
			TenantUid: types.StringValue(sdc.TenantUid),

			SdcPublicKey: PublicKeyModel{
				EncodedKey: types.StringValue(sdc.PublicKey.EncodedKey),
				Version:    types.Int64Value(sdc.PublicKey.Version),
				KeyId:      types.StringValue(sdc.PublicKey.KeyId),
			},
		}

		sdcModels = append(sdcModels, sdcModel)
	}
	configData.Sdcs = sdcModels

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &configData)...)
}
