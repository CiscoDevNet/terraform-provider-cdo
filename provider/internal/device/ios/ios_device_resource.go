package ios

import (
	"context"
	"fmt"
	"github.com/cisco-lockhart/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strconv"

	cdoClient "github.com/cisco-lockhart/go-client"
	"github.com/cisco-lockhart/go-client/device/ios"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	SdcUid  types.String `tfsdk:"sdc_uid"`
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
			"sdc_uid": schema.StringAttribute{
				MarkdownDescription: "The SDC UID that will be used to communicate with the device",
				Optional:            true,
				Computed:            true,
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
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password used to authenticate with the device",
				Required:            true,
				Sensitive:           true,
			},
			"ignore_certificate": schema.BoolAttribute{
				MarkdownDescription: "Whether to ignore certificate validation",
				Optional:            true,
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

	var stateData *IosDeviceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read ios
	readInp := ios.ReadInput{
		Uid: stateData.ID.ValueString(),
	}
	readOutp, err := r.client.ReadIos(ctx, readInp)
	if err != nil {
		resp.Diagnostics.AddError("unable to read IOS Device", err.Error())
		return
	}

	port, err := strconv.ParseInt(readOutp.Port, 10, 16)
	if err != nil {
		resp.Diagnostics.AddError("unable to read IOS Device", err.Error())
		return
	}
	stateData.Port = types.Int64Value(port)

	stateData.ID = types.StringValue(readOutp.Uid)
	stateData.SdcType = types.StringValue(readOutp.LarType)
	stateData.SdcUid = types.StringValue(readOutp.LarUid)
	stateData.Name = types.StringValue(readOutp.Name)
	stateData.Ipv4 = types.StringValue(readOutp.Ipv4)
	stateData.Host = types.StringValue(readOutp.Host)
	stateData.IgnoreCertifcate = types.BoolValue(readOutp.IgnoreCertifcate)

	// Fix: where to find them? We need them for import statement
	// stateData.Username = types.StringNull()
	// stateData.Password = types.StringNull()

	tflog.Trace(ctx, "done read IOS device resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *IosDeviceResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	tflog.Trace(ctx, "create IOS device resource")

	var planData IosDeviceResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	createInp := ios.NewCreateRequestInput(
		planData.Name.ValueString(),
		planData.SdcUid.ValueString(),
		planData.SdcType.ValueString(),
		planData.Ipv4.ValueString(),
		planData.Username.ValueString(),
		planData.Password.ValueString(),
		planData.IgnoreCertifcate.ValueBool(),
	)

	createOutp, err := r.client.CreateIos(ctx, *createInp)
	if err != nil {
		res.Diagnostics.AddError("failed to create IOS", err.Error())
		return
	}

	planData.ID = types.StringValue(createOutp.Uid)
	planData.SdcType = types.StringValue(createOutp.LarType)
	planData.SdcUid = types.StringValue(createOutp.LarUid)
	planData.Name = types.StringValue(createOutp.Name)
	planData.Host = types.StringValue(createOutp.Host)

	port, err := strconv.ParseInt(createOutp.Port, 10, 16)
	if err != nil {
		res.Diagnostics.AddError("failed to parse IOS port", err.Error())
	}
	planData.Port = types.Int64Value(port)

	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
}

func (r *IosDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {

	tflog.Trace(ctx, "update IOS device resource")

	var planData *IosDeviceResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	var stateData *IosDeviceResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	updateInp := ios.NewUpdateInput(
		stateData.ID.ValueString(),
		planData.Name.ValueString(),
	)

	updateOutp, err := r.client.UpdateIos(ctx, *updateInp)
	if err != nil {
		res.Diagnostics.AddError("failed to update IOS device", err.Error())
		return
	}

	stateData.Name = types.StringValue(updateOutp.Name)

	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
}

func (r *IosDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {

	tflog.Trace(ctx, "delete IOS device resource")

	var stateData IosDeviceResourceModel

	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	deleteInp := ios.NewDeleteInput(stateData.ID.ValueString())
	_, err := r.client.DeleteIos(ctx, *deleteInp)
	if err != nil {
		res.Diagnostics.AddError("failed to delete IOS device", err.Error())
		return
	}

}

func (r *IosDeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}
