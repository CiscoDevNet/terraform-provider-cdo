// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sdc

import (
	"context"
	"fmt"

	cdoClient "github.com/cisco-lockhart/go-client"
	"github.com/cisco-lockhart/go-client/device/sdc"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Used in provider.go to include this data source.
func NewSdcDataSource() datasource.DataSource {
	return &SdcDataSource{}
}

// The data source object consumed by terraform.
type SdcDataSource struct {
	client *cdoClient.Client
}

/////
// Model classes: define mapping from Go types to Terraform types.
/////

type SdcDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Sdcs []SdcModel   `tfsdk:"sdcs"`
}

type SdcModel struct {
	Uid          types.String      `tfsdk:"uid"`
	Name         types.String      `tfsdk:"name"`
	TenantUid    types.String      `tfsdk:"tenant_uid"`
	SdcPublicKey SdcPublicKeyModel `tfsdk:"sdc_public_key"`
}

type SdcPublicKeyModel struct {
	EncodedKey types.String `tfsdk:"encoded_key"`
	Version    types.Int64  `tfsdk:"version"`
	KeyId      types.String `tfsdk:"key_id"`
}

// define the name for this data source.
func (d *SdcDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc_device"
}

// define the terraform schema for this data source.
// it needs to be consistent with the Model classes' `tfsdk:"xxx"` above.
func (d *SdcDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
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

// initial configuration for CRUD operations, here we set the cdo go client to use.
// the go client is configured in the provider, and it is available here, so we just set it directly.
func (d *SdcDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// define the read operation.
// this function is responsible for:
// 1. mapping the cdo go client's data to the model classes we defined above.
// 2. report any error using `resp.Diagnostics`.
// 3. set model classes to the terraform state.
func (d *SdcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var configData SdcDataSourceModel

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

	var sdcModels []SdcModel
	for _, sdc := range *res {
		sdcModel := SdcModel{
			Uid:       types.StringValue(sdc.Uid),
			Name:      types.StringValue(sdc.Name),
			TenantUid: types.StringValue(sdc.TenantUid),

			SdcPublicKey: SdcPublicKeyModel{
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
