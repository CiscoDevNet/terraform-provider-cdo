package sec

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *cdoClient.Client
}

type ResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	CdoBootstrapData types.String `tfsdk:"cdo_bootstrap_data"`
	SecBootstrapData types.String `tfsdk:"sec_bootstrap_data"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sec"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides an SEC connector resource. This allows SEC to be onboarded, updated, and deleted on CDO.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the device. This is a UUID and is automatically generated when the device is created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A generated name for the Secure Event Connector (SEC). This name is unique.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cdo_bootstrap_data": schema.StringAttribute{
				MarkdownDescription: "CDO bootstrap data.",
				Computed:            true,
				Sensitive:           true, // bootstrap data contains user api token
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sec_bootstrap_data": schema.StringAttribute{
				MarkdownDescription: "SEC bootstrap data.",
				Computed:            true,
				Sensitive:           true, // bootstrap data contains user api token
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Trace(ctx, "read SEC resource")

	// 1. read state data
	var stateData ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	if err := Read(ctx, r, &stateData); err != nil {
		resp.Diagnostics.AddError("failed to read SEC resource", err.Error())
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	tflog.Trace(ctx, "create Sec resource")

	// 1. read plan data into planData
	var planData ResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. use plan data to create device and fill up rest of the model
	if err := Create(ctx, r, &planData); err != nil {
		res.Diagnostics.AddError("failed to create Sec resource", err.Error())
		return
	}

	// 3. set state using filled model
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {

	tflog.Trace(ctx, "update Sec resource")

	// 1. read plan data
	var planData ResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. read state data
	var stateData ResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 3. do update
	if err := Update(ctx, r, &planData, &stateData); err != nil {
		res.Diagnostics.AddError("failed to update Sec resource", err.Error())
	}

	// 4. set resulting state
	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {

	tflog.Trace(ctx, "delete Sec resource")

	var stateData ResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	if err := Delete(ctx, r, &stateData); err != nil {
		res.Diagnostics.AddError("failed to delete Sec resource", err.Error())
	}
}
