package ios

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &IosDeviceResource{}
var _ resource.ResourceWithImportState = &IosDeviceResource{}

func NewIosDeviceResource() resource.Resource {
	return &IosDeviceResource{}
}

type IosDeviceResource struct {
	client *cdoClient.Client
}

type IosDeviceResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ConnectorName types.String `tfsdk:"connector_name"`
	Name          types.String `tfsdk:"name"`
	Ipv4          types.String `tfsdk:"socket_address"`
	Host          types.String `tfsdk:"host"`
	Port          types.Int64  `tfsdk:"port"`
	Labels        types.Set    `tfsdk:"labels"`

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`

	IgnoreCertificate types.Bool `tfsdk:"ignore_certificate"`
}

func (r *IosDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ios_device"
}

func (r *IosDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides an iOS device resource. This allows iOS devices to be onboarded, updated, and deleted on CDO.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the device. This is a UUID and is automatically generated when the device is created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A human-readable name for the device.",
				Required:            true,
			},
			"connector_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Secure Device Connector (SDC) that will be used to communicate with the device. This value is not required if the connector type selected is Cloud Device Gateway (CDG).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"socket_address": schema.StringAttribute{
				MarkdownDescription: "The address of the device to onboard, specified in the format `host:port`.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validators.ValidateSocketAddress(),
				},
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port used to connect to the device.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host used to connect to the device.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username used to authenticate with the device.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password used to authenticate with the device.",
				Required:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ignore_certificate": schema.BoolAttribute{
				MarkdownDescription: "Set this attribute to true if you do not want CDO to validate the certificate of this device before onboarding.",
				Required:            true,
			},
			"labels": schema.SetAttribute{
				MarkdownDescription: "Set a set of labels to identify the device as part of a group. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Optional:            true,
				ElementType:         types.StringType,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})), // default to empty list
			},
		},
	}
}

func (r *IosDeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdoClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cdoClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IosDeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Trace(ctx, "read IOS device resource")

	// 1. read state data
	var stateData IosDeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	if err := Read(ctx, r, &stateData); err != nil {
		resp.Diagnostics.AddError("failed to read IOS device", err.Error())
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *IosDeviceResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	tflog.Trace(ctx, "create IOS device resource")

	// 1. read plan data into planData
	var planData IosDeviceResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. use plan data to create device and fill up rest of the model
	if err := Create(ctx, r, &planData); err != nil {
		res.Diagnostics.AddError("failed to create IOS device", err.Error())
		return
	}

	// 3. set state using filled model
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
}

func (r *IosDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {

	tflog.Trace(ctx, "update IOS device resource")

	// 1. read plan data
	var planData IosDeviceResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. read state data
	var stateData IosDeviceResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 3. do update
	if err := Update(ctx, r, &planData, &stateData); err != nil {
		res.Diagnostics.AddError("failed to update IOS device", err.Error())
	}

	// 4. set resulting state
	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
}

func (r *IosDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {

	tflog.Trace(ctx, "delete IOS device resource")

	var stateData IosDeviceResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	if err := Delete(ctx, r, &stateData); err != nil {
		res.Diagnostics.AddError("failed to delete IOS device", err.Error())
	}
}

func (r *IosDeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func (r *IosDeviceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, res *resource.ModifyPlanResponse) {
	if !req.State.Raw.IsNull() {
		// this is an update
		var stateData *IosDeviceResourceModel
		res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		if res.Diagnostics.HasError() {
			return
		}

		var planData *IosDeviceResourceModel
		res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		if res.Diagnostics.HasError() {
			return
		}

		if planData != nil && stateData != nil && strings.EqualFold(planData.Ipv4.ValueString(), stateData.Ipv4.ValueString()) {
			tflog.Debug(ctx, "There is no change in the IPv4; remove host and port diffs")
			planData.Host = stateData.Host
			planData.Port = stateData.Port
		}

		res.Diagnostics.Append(res.Plan.Set(ctx, &planData)...)
	}
}
