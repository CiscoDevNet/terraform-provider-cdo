// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ios

import (
	"context"
	"fmt"
	"strconv"

	cdoClient "github.com/cisco-lockhart/go-client"
	"github.com/cisco-lockhart/go-client/device/ios"
	"github.com/cisco-lockhart/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = NewIosDataSource()

// Used in provider.go to include this data source.
func NewIosDataSource() datasource.DataSource {
	return &IosDataSource{}
}

// The data source object consumed by terraform.
type IosDataSource struct {
	client *cdoClient.Client
}

/////
// Model classes: define mapping from Go types to Terraform types.
/////

type IosDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	SdcType          types.String `tfsdk:"sdc_type"`
	SdcName          types.String `tfsdk:"sdc_name"`
	Name             types.String `tfsdk:"name"`
	Ipv4             types.String `tfsdk:"ipv4"`
	Host             types.String `tfsdk:"host"`
	Port             types.Int64  `tfsdk:"port"`
	IgnoreCertifcate types.Bool   `tfsdk:"ignore_certificate"`
}

// define the name for this data source.
func (d *IosDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ios_device"
}

// define the terraform schema for this data source.
// it needs to be consistent with the Model classes' `tfsdk:"xxx"` above.
func (d *IosDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "IOS data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Uid used to represent the device",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name assigned to the device",
				Computed:            true,
			},
			"sdc_name": schema.StringAttribute{
				MarkdownDescription: "The SDC name that will be used to communicate with the device",
				Computed:            true,
			},
			"sdc_type": schema.StringAttribute{
				MarkdownDescription: "The type of SDC that will be used to communicate with the device (Valid values: [CDG, SDC])",
				Computed:            true,
				Validators: []validator.String{
					validators.OneOf("CDG", "SDC"),
				},
			},
			"ipv4": schema.StringAttribute{
				MarkdownDescription: "The ipv4 address of the device",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port used to connect to the device",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host used to connect to the device",
				Computed:            true,
			},
			"ignore_certificate": schema.BoolAttribute{
				MarkdownDescription: "Whether to ignore certificate validation",
				Computed:            true,
			},
		},
	}
}

// initial configuration for CRUD operations, here we set the cdo go client to use.
// the go client is configured in the provider, and it is available here, so we just set it directly.
func (d *IosDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *IosDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Trace(ctx, "read IOS device data source")

	var configData *IosDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read ios
	readInp := ios.ReadInput{
		Uid: configData.ID.ValueString(),
	}
	readOutp, err := d.client.ReadIos(ctx, readInp)
	if err != nil {
		resp.Diagnostics.AddError("unable to read IOS Device", err.Error())
		return
	}

	port, err := strconv.ParseInt(readOutp.Port, 10, 16)
	if err != nil {
		resp.Diagnostics.AddError("unable to read IOS Device", err.Error())
		return
	}
	configData.Port = types.Int64Value(port)

	configData.ID = types.StringValue(readOutp.Uid)
	configData.SdcType = types.StringValue(readOutp.LarType)
	configData.Name = types.StringValue(readOutp.Name)
	configData.Ipv4 = types.StringValue(readOutp.Ipv4)
	configData.Host = types.StringValue(readOutp.Host)
	configData.IgnoreCertifcate = types.BoolValue(readOutp.IgnoreCertifcate)

	tflog.Trace(ctx, "done read IOS device data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &configData)...)
}
