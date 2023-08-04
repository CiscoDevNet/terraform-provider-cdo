package sdc

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevNet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &SdcResource{}
var _ resource.ResourceWithImportState = &SdcResource{}

func NewSdcResource() resource.Resource {
	return &SdcResource{}
}

type SdcResource struct {
	client *cdoClient.Client
}

type SdcResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	BootstrapData types.String `tfsdk:"bootstrap_data"`
}

func (r *SdcResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

func (r *SdcResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides an SDC connector resource. This allows SDC to be onboarded, updated, and deleted on CDO.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the device. This is a UUID and will be automatically generated when the device is created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A human-readable name for the Secure Device Connector (SDC). This should be unique among SDCs",
				Required:            true,
			},
			"bootstrap_data": schema.StringAttribute{
				MarkdownDescription: "SDC bootstrap data",
				Computed:            true,
				// TODO: uncomment to make it sensitive
				// Sensitive:           true, // bootstrap data contais user api token
			},
		},
	}
}

func (r *SdcResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "read SDC resource")

	// 1. read terraform plan data into the model
	var stateData SdcResourceModel
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

func (r *SdcResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	tflog.Trace(ctx, "create SDC resource")

	// 1. read terraform plan data into model
	var planData SdcResourceModel
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

func (r *SdcResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	tflog.Trace(ctx, "update SDC resource")

	// 1. read plan and state data from terraform
	var planData SdcResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}
	var stateData SdcResourceModel
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

func (r *SdcResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	tflog.Trace(ctx, "delete SDC resource")

	// 1. read state data from terraform state
	var stateData SdcResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. delete the resource
	if err := Delete(ctx, r, &stateData); err != nil {
		res.Diagnostics.AddError("failed to delete SDC resource", err.Error())
	}
}

func (r *SdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}
