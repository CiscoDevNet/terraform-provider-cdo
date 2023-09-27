// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ios

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	ID                types.String   `tfsdk:"id"`
	Name              types.String   `tfsdk:"name"`
	ConnectorName     types.String   `tfsdk:"connector_name"`
	SocketAddress     types.String   `tfsdk:"socket_address"`
	Host              types.String   `tfsdk:"host"`
	Port              types.Int64    `tfsdk:"port"`
	IgnoreCertificate types.Bool     `tfsdk:"ignore_certificate"`
	Labels            []types.String `tfsdk:"labels"`
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
				MarkdownDescription: "Universally unique identifier of the device.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-readable name of the device. This is the name displayed on the CDO Inventory page. Device names are unique across a CDO tenant.",
				Required:            true,
			},
			"connector_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Secure Device Connector (SDC) that is used by CDO to communicate with the device.",
				Computed:            true,
			},
			"socket_address": schema.StringAttribute{
				MarkdownDescription: "The address of the device to onboard, specified in the format `host:port`.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port used to connect to the device.",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host used to connect to the device.",
				Computed:            true,
			},
			"ignore_certificate": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether certificates were ignored when onboarding this device.",
				Computed:            true,
			},
			"labels": schema.ListAttribute{
				MarkdownDescription: "Set a list of labels to identify the device as part of a group. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.UniqueValues(),
				},
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
	readInp := device.ReadByNameAndTypeInput{
		Name:       configData.Name.ValueString(),
		DeviceType: "IOS",
	}
	readOutp, err := d.client.ReadDeviceByName(ctx, readInp)
	if err != nil {
		resp.Diagnostics.AddError("unable to find IOS Device", err.Error())
		return
	}

	port, err := strconv.ParseInt(readOutp.Port, 10, 16)
	if err != nil {
		resp.Diagnostics.AddError("unable to find IOS Device", err.Error())
		return
	}
	configData.Port = types.Int64Value(port)

	configData.ID = types.StringValue(readOutp.Uid)
	configData.Name = types.StringValue(readOutp.Name)
	configData.SocketAddress = types.StringValue(readOutp.SocketAddress)
	configData.Host = types.StringValue(readOutp.Host)
	configData.IgnoreCertificate = types.BoolValue(readOutp.IgnoreCertificate)
	configData.Labels = util.GoStringSliceToTFStringList(readOutp.Tags.Labels)

	tflog.Trace(ctx, "done read IOS device data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &configData)...)
}
