// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package asa

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = NewAsaDataSource()

// Used in provider.go to include this data source.
func NewAsaDataSource() datasource.DataSource {
	return &AsaDataSource{}
}

// The data source object consumed by terraform.
type AsaDataSource struct {
	client *cdoClient.Client
}

/////
// Model classes: define mapping from Go types to Terraform types.
/////

type AsaDataSourceModel struct {
	ID                types.String   `tfsdk:"id"`
	ConnectorType     types.String   `tfsdk:"connector_type"`
	SdcName           types.String   `tfsdk:"sdc_name"`
	Name              types.String   `tfsdk:"name"`
	Ipv4              types.String   `tfsdk:"socket_address"`
	Host              types.String   `tfsdk:"host"`
	Port              types.Int64    `tfsdk:"port"`
	IgnoreCertificate types.Bool     `tfsdk:"ignore_certificate"`
	Labels            []types.String `tfsdk:"labels"`
	GroupedLabels     types.Map      `tfsdk:"grouped_labels"`
}

// define the name for this data source.
func (d *AsaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asa_device"
}

// define the terraform schema for this data source.
// it needs to be consistent with the Model classes' `tfsdk:"xxx"` above.
func (d *AsaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "ASA data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the device.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The human-readable name of the device. This is the name displayed on the CDO Inventory page. Device names are unique across a CDO tenant.",
				Required:            true,
			},
			"sdc_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Secure Device Connector (SDC) that is used by CDO to communicate with the device. This value will be empty if the device was onboarded using a Cloud Connector (CDG).",
				Computed:            true,
			},
			"connector_type": schema.StringAttribute{
				MarkdownDescription: "The type of the connector that is used to communicate with the device. CDO can communicate with your device using either a Cloud Connector (CDG) or a Secure Device Connector (SDC); see [the CDO documentation](https://docs.defenseorchestrator.com/c-connect-cisco-defense-orchestratortor-the-secure-device-connector.html) to learn more (Valid values: [CDG, SDC]).",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("CDG", "SDC"),
				},
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
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "The labels applied to the device. Labels are used to group devices in CDO. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Validators: []validator.List{
					listvalidator.UniqueValues(),
				},
			},
			"grouped_labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				Computed:            true,
				MarkdownDescription: "The grouped labels applied to the device. Labels are used to group devices in CDO. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
			},
		},
	}
}

// initial configuration for CRUD operations, here we set the cdo go client to use.
// the go client is configured in the provider, and it is available here, so we just set it directly.
func (d *AsaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *AsaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Trace(ctx, "read ASA device data source")

	var configData *AsaDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read asa
	readInp := device.ReadByNameAndTypeInput{
		Name:       configData.Name.ValueString(),
		DeviceType: "ASA",
	}
	readOutp, err := d.client.ReadDeviceByName(ctx, readInp)
	if err != nil {
		resp.Diagnostics.AddError("unable to find ASA Device", err.Error())
		return
	}

	port, err := strconv.ParseInt(readOutp.Port, 10, 16)
	if err != nil {
		resp.Diagnostics.AddError("unable to find ASA Device", err.Error())
		return
	}
	configData.Port = types.Int64Value(port)

	configData.ID = types.StringValue(readOutp.Uid)
	configData.ConnectorType = types.StringValue(readOutp.ConnectorType)
	configData.Name = types.StringValue(readOutp.Name)
	configData.Ipv4 = types.StringValue(readOutp.SocketAddress)
	configData.Host = types.StringValue(readOutp.Host)
	configData.IgnoreCertificate = types.BoolValue(readOutp.IgnoreCertificate)
	configData.Labels = util.GoStringSliceToTFStringList(readOutp.Tags.UngroupedTags())
	configData.GroupedLabels = util.GoMapToStringSetTFMap(readOutp.Tags.GroupedTags())

	tflog.Trace(ctx, "done read ASA device data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &configData)...)
}
