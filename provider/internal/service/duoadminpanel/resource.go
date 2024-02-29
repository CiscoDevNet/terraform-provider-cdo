package duoadminpanel

import (
	"context"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	IntegrationKey types.String `tfsdk:"integration_key"`
	SecretKey      types.String `tfsdk:"secret_key"`
	Host           types.String `tfsdk:"host"`
	Labels         types.Set    `tfsdk:"labels"`
	GroupedLabels  types.Map    `tfsdk:"grouped_labels"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_duo_admin_panel"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides an Duo Admin Panel resource. This allows Duo Admin Panels to be onboarded, updated, and deleted on CDO.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the Duo Admin Panel. This is a UUID and is automatically generated when the Duo Admin Panel is created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A human-readable name for the Duo Admin Panel.",
				Required:            true,
			},
			"integration_key": schema.StringAttribute{
				MarkdownDescription: "The integration key for an Admin API application in the Duo Admin Panel. Refer to the CDO documentation for details on how to create an Admin API application to onboard your Duo Admin Panel in CDO.",
				Required:            true,
				Sensitive:           true,
			},
			"secret_key": schema.StringAttribute{
				MarkdownDescription: "The secret key for an Admin API application in the Duo Admin Panel. Refer to the CDO documentation for details on how to create an Admin API application to onboard your Duo Admin Panel in CDO.",
				Required:            true,
				Sensitive:           true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The API hostname for an Admin API application in the Duo Admin Panel. Refer to the CDO documentation for details on how to create an Admin API application to onboard your Duo Admin Panel in CDO.",
				Required:            true,
			},
			"labels": schema.SetAttribute{
				MarkdownDescription: "Specify a set of labels to identify the Duo Admin Panel as part of a group. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})), // default to empty list
			},
			"grouped_labels": schema.MapAttribute{
				MarkdownDescription: "Specify a set of grouped labels to identify the Duo Admin Panel as part of a group. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Optional:            true,
				Computed:            true,
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				Default: mapdefault.StaticValue(types.MapValueMust(types.SetType{ElemType: types.StringType}, map[string]attr.Value{})), // default to empty map
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

	tflog.Trace(ctx, "read Duo Admin Panel resource")

	// 1. read state data
	var stateData ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	if err := Read(ctx, r, &stateData); err != nil {
		if util.Is404Error(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read Duo Admin Panel device", err.Error())
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	tflog.Trace(ctx, "create Duo Admin Panel resource")

	// 1. read plan data into planData
	var planData ResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. use plan data to create device and fill up rest of the model
	if err := Create(ctx, r, &planData); err != nil {
		res.Diagnostics.AddError("failed to create Duo Admin Panel resource", err.Error())
		return
	}

	// 3. set state using filled model
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {

	tflog.Trace(ctx, "update Duo Admin Panel resource")

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
		res.Diagnostics.AddError("failed to update Duo Admin Panel resource", err.Error())
	}

	// 4. set resulting state
	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {

	tflog.Trace(ctx, "delete Duo Admin Panel resource")

	var stateData ResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	if err := Delete(ctx, r, &stateData); err != nil {
		res.Diagnostics.AddError("failed to delete Duo Admin Panel resource", err.Error())
	}
}
