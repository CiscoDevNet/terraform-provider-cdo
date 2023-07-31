package ios

import (
	"context"
	"fmt"

	"github.com/cisco-lockhart/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	cdoClient "github.com/cisco-lockhart/go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
	ID      types.String `tfsdk:"id"`
	SdcType types.String `tfsdk:"sdc_type"`
	SdcName types.String `tfsdk:"sdc_name"`
	Name    types.String `tfsdk:"name"`
	Ipv4    types.String `tfsdk:"ipv4"`
	Host    types.String `tfsdk:"host"`
	Port    types.Int64  `tfsdk:"port"`

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`

	IgnoreCertifcate types.Bool `tfsdk:"ignore_certificate"`
}

func (r *IosDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ios_device"
}

func (r *IosDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "IOS Device resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Uid used to represent the device",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name to assign the device",
				Required:            true,
			},
			"sdc_name": schema.StringAttribute{
				MarkdownDescription: "The SDC name that will be used to communicate with the device",
				Required:            true,
			},
			"sdc_type": schema.StringAttribute{
				MarkdownDescription: "The type of SDC that will be used to communicate with the device (Valid values: [CDG, SDC])",
				Required:            true,
				Validators: []validator.String{
					validators.OneOf("CDG", "SDC"),
				},
			},
			"ipv4": schema.StringAttribute{
				MarkdownDescription: "The ipv4 address of the device",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port used to connect to the device",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host used to connect to the device",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username used to authenticate with the device",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password used to authenticate with the device",
				Required:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ignore_certificate": schema.BoolAttribute{
				MarkdownDescription: "Whether to ignore certificate validation",
				Computed:            true,
				Default:             booldefault.StaticBool(false),
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
