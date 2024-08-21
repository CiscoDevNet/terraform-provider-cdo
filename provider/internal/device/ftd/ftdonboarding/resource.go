package ftdonboarding

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
	Id     types.String `tfsdk:"id"`
	FtdUid types.String `tfsdk:"ftd_uid"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ftd_device_onboarding"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource is meant to be in conjunction with the `cdo_ftd_device` resource to complete the onboarding process of an FTD to a cdFMC. " +
			"The `cdo_ftd_device` creates an FTD device on CDO and generates a command with the registration key that should be pasted into the FTD device's CLI over SSH (see **step 10** [here](https://docs.defenseorchestrator.com/c_onboard-an-ftd.html#!t-onboard-an-ftd-device-with-regkey.html)). " +
			"This resource waits for you to finish pasting the registration command, and onboards the FTD to CDO. " +
			"If you are spinning up an FTDv using Terraform, you can pass the output of the `cdo_ftd_device` through to the FTDv. " +
			"If you are using a manually deployed FTDv or a physical FTD, you cannot add this resource to your Terraform code until after you have applied the `cdo_ftd_device` resource, retrieved the `generated_command` from the resource, and pasted it into the FTD device's CLI. " +
			"This resource will time out if the registration command is not applied on the FTD CLI within 3 minutes of it starting to poll.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of this FTD onboarding resource, it is the registration key of the FTD.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ftd_uid": schema.StringAttribute{
				MarkdownDescription: "The ID of the FTD to add to the cdFMC. This value is returned by the `id` attribute of the `cdo_ftd_device` resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
	tflog.Trace(ctx, "read FTD onboarding resource")

	// 1. read terraform plan data into the model
	var stateData ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	if err := Read(ctx, r, &stateData); err != nil {
		resp.Diagnostics.AddError("failed to read FTD onboarding resource", err.Error())
		return
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	tflog.Trace(ctx, "read FTD onboarding resource done")
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	tflog.Trace(ctx, "create FTD onboarding resource")

	// 1. read terraform plan data into model
	var planData ResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. create resource & fill model data
	if err := Create(ctx, r, &planData); err != nil {
		res.Diagnostics.AddError("failed to create FTD onboarding resource", err.Error())
		return
	}

	// 3. fill terraform state using model data
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
	tflog.Trace(ctx, "create FTD onboarding resource done")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	tflog.Trace(ctx, "update FTD onboarding resource")

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
		res.Diagnostics.AddError("failed to update FTD onboarding resource", err.Error())
		return
	}

	// 3. update terraform state with updated state data
	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
	tflog.Trace(ctx, "update FTD onboarding resource done")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	tflog.Trace(ctx, "delete FTD onboarding resource")

	// 1. read state data from terraform state
	var stateData ResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. delete the resource
	if err := Delete(ctx, r, &stateData); err != nil {
		res.Diagnostics.AddError("failed to delete FTD onboarding resource", err.Error())
	}
}
