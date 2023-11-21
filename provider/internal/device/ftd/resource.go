package ftd

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	AccessPolicyName types.String `tfsdk:"access_policy_name"`
	PerformanceTier  types.String `tfsdk:"performance_tier"`
	Virtual          types.Bool   `tfsdk:"virtual"`
	Licenses         types.Set    `tfsdk:"licenses"`
	Labels           types.Set    `tfsdk:"labels"`

	AccessPolicyUid  types.String `tfsdk:"access_policy_id"`
	GeneratedCommand types.String `tfsdk:"generated_command"`
	Hostname         types.String `tfsdk:"hostname"`
	NatId            types.String `tfsdk:"nat_id"`
	RegKey           types.String `tfsdk:"reg_key"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ftd_device" // TODO: _cloud_ftd_device ?
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Firewall Threat Defense device resource. Use this to onboard, update, and delete FTDs from CDO. This resource does not complete the onboarding of an FTD into CDO and cdFMC. " +
			"It creates a FTD device entry in the CDO Inventory, and generates a registration command (see the `generated_command` attribute) that needs to be pasted into the FTD CLI (see **step 10** [here](https://docs.defenseorchestrator.com/c_onboard-an-ftd.html#!t-onboard-an-ftd-device-with-regkey.html)). " +
			"To finish adding the FTD device to CDO and cdFMC, use the `cdo_ftd_device_onboarding` resource after you have applied this resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the device. This is a UUID and is automatically generated when the device is created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A human-readable name for the Firewall Threat Defense (FTD). This name must be unique.",
				Required:            true,
			},
			"access_policy_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Cloud-Delivered FMC (cdFMC) access policy that will be used by the FTD.",
				Required:            true,
				// TODO: make this optional, and use default access policy when not given
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"performance_tier": schema.StringAttribute{
				MarkdownDescription: "The performance tier of the virtual FTD, if virtual is set to false, this field is ignored as performance tiers are not applicable to physical FTD devices. Allowed values are: [\"FTDv5\", \"FTDv10\", \"FTDv20\", \"FTDv30\", \"FTDv50\", \"FTDv100\", \"FTDv\"].",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(tier.AllAsString...),
				},
			},
			"virtual": schema.BoolAttribute{
				MarkdownDescription: "This determines if this FTD is virtual. If false, performance_tier is ignored as performance tiers are not applicable to physical FTD devices.",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"licenses": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Comma-separated list of licenses to apply to this FTD. You must enable at least the \"BASE\" license. Allowed values are: [\"BASE\", \"CARRIER\", \"THREAT\", \"MALWARE\", \"URLFilter\",].",
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(stringvalidator.OneOf(license.AllAsString...)),
					validators.SetValueStringsAtLeast(stringvalidator.OneOf(string(license.Base), string(license.Essentials))),
				},
			},
			"labels": schema.SetAttribute{
				MarkdownDescription: "Specify a set of labels to identify the device as part of a group. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Optional:            true,
				ElementType:         types.StringType,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})), // default to empty set
			},
			"generated_command": schema.StringAttribute{
				MarkdownDescription: "The command to run in the FTD CLI to register it with the cloud-delivered FMC (cdFMC).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_policy_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the cloud-delivered FMC (cdFMC) access policy applied to this FTD.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The Hostname of the cloud-delivered FMC (cdFMC) manages this FTD.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nat_id": schema.StringAttribute{
				MarkdownDescription: "The Network Address Translation (NAT) ID of this FTD.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"reg_key": schema.StringAttribute{
				MarkdownDescription: "The Registration Key of this FTD.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
	tflog.Trace(ctx, "read FTD resource")

	// 1. read terraform plan data into the model
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
		resp.Diagnostics.AddError("failed to read FTD resource", err.Error())
		return
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	tflog.Trace(ctx, "read FTD resource done")
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	tflog.Trace(ctx, "create FTD resource")

	// 1. read terraform plan data into model
	var planData ResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. create resource & fill model data
	if err := Create(ctx, r, &planData); err != nil {
		res.Diagnostics.AddError("failed to create FTD resource", err.Error())
		return
	}

	// 3. fill terraform state using model data
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
	tflog.Trace(ctx, "create FTD resource done")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	tflog.Trace(ctx, "update FTD resource")

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
		res.Diagnostics.AddError("failed to update FTD resource", err.Error())
		return
	}

	// 3. update terraform state with updated state data
	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
	tflog.Trace(ctx, "update FTD resource done")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	tflog.Trace(ctx, "delete FTD resource")

	// 1. read state data from terraform state
	var stateData ResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. delete the resource
	if err := Delete(ctx, r, &stateData); err != nil {
		res.Diagnostics.AddError("failed to delete FTD resource", err.Error())
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}
