package connectoronboarding

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *cdoClient.Client
}

type ResourceModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc_onboarding"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this resource to wait for an SDC to finish onboarding. When an SDC is onboarded, either manually or using the CDO Terraform Modules for [AWS](linkhere) and [vSphere](linkhere), it can take a few minutes before the SDC is active and capable of proxying communications between CDO and the device. This resource allows you to wait until this is done.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of this SDC onboarding resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Specify the name of the SDC.",
				Required:            true,
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	tflog.Trace(ctx, "read SDC resource")

	// 1. read terraform plan data into the model
	var stateData ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	if err := Read(ctx, r, &stateData); err != nil {
		resp.Diagnostics.AddError("failed to read SDC resource", err.Error())
		return
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	tflog.Trace(ctx, "read SDC resource done")
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	tflog.Trace(ctx, "create SDC resource")

	// 1. read terraform plan data into model
	var planData ResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. create resource & fill model data
	if err := Create(ctx, r, &planData); err != nil {
		res.Diagnostics.AddError("failed to create SDC resource", err.Error())
		return
	}

	// 3. fill terraform state using model data
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
	tflog.Trace(ctx, "create SDC resource done")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	tflog.Trace(ctx, "update SDC resource")

	// 1. read plan and state data from terraform
	var planData ResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}
	var stateData ResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. update resource & state data
	if err := Update(ctx, r, &planData, &stateData); err != nil {
		res.Diagnostics.AddError("failed to update SDC resource", err.Error())
		return
	}

	// 3. update terraform state with updated state data
	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
	tflog.Trace(ctx, "update SDC resource done")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	tflog.Trace(ctx, "delete SDC resource")

	// 1. read state data from terraform state
	var stateData ResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. delete the resource
	if err := Delete(ctx, r, &stateData); err != nil {
		res.Diagnostics.AddError("failed to delete SDC resource", err.Error())
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}
