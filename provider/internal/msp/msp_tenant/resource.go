package msp_tenant

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/tenants"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewTenantResource() resource.Resource { return &TenantResource{} }

type TenantResource struct {
	client *cdoClient.Client
}

type TenantResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DisplayName   types.String `tfsdk:"display_name"`
	GeneratedName types.String `tfsdk:"generated_name"`
	Region        types.String `tfsdk:"region"`
}

func (*TenantResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_msp_managed_tenant"
}

func (*TenantResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Provides an MSP managed tenant resource. This allows MSP managed tenants to be created.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the tenant",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the tenant",
				Validators: []validator.String{
					validators.NewCdoTenantValidator(),
				},
				Required: true,
				PlanModifiers: []planmodifier.String{
					PreventUpdatePlanModifier{}, // Prevent updates to name
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the tenant. If no display name is specified, the display name will be set to the tenant name.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					PreventUpdatePlanModifier{}, // Prevent updates to name
				},
			},
			"generated_name": schema.StringAttribute{
				MarkdownDescription: "Actual name of the tenant returned by the API. This auto-generated name will differ from the name entered by the customer.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "CDO region in which the tenant is created. This is the same region as the region of the MSP portal.",
				Computed:            true,
			},
		},
	}
}

func (resource *TenantResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdoClient.Client)

	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cdoClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	resource.client = client
}

func (t *TenantResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating a CDO tenant")

	// 1. Read plan data into planData
	var planData TenantResourceModel

	response.Diagnostics.Append(request.Plan.Get(ctx, &planData)...)
	tflog.Debug(ctx, fmt.Sprintf("Diagnostics: %v", response.Diagnostics))
	tflog.Debug(ctx, fmt.Sprintf("lavda: %v", planData.Name))

	if response.Diagnostics.HasError() {
		return
	}

	// 2. use plan data to create tenant and fill up rest of the model
	createOut, err := t.client.CreateTenantUsingMspPortal(ctx, tenants.MspCreateTenantInput{
		Name:        planData.Name.ValueString(),
		DisplayName: planData.DisplayName.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError("failed to create CDO Tenant", err.Error())
		return
	}

	planData.Id = types.StringValue(createOut.Uid)
	planData.Name = types.StringValue(planData.Name.ValueString())
	planData.DisplayName = types.StringValue(createOut.DisplayName)
	planData.GeneratedName = types.StringValue(createOut.Name)
	planData.Region = types.StringValue(createOut.Region)
	response.Diagnostics.Append(response.State.Set(ctx, &planData)...)
}

func (t *TenantResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading a CDO tenant")
	var stateData *TenantResourceModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)
	if response.Diagnostics.HasError() {
		return
	}

	// read Tenant
	readInp := tenants.ReadByUidInput{
		Uid: stateData.Id.ValueString(),
	}
	tenantReadOutp, err := t.client.ReadMspManagedTenant(ctx, readInp)

	if err != nil {
		if util.Is404Error(err) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("unable to read tenant", err.Error())
		return
	}

	stateData.Id = types.StringValue(tenantReadOutp.Uid)
	stateData.Name = types.StringValue(stateData.Name.ValueString())
	stateData.GeneratedName = types.StringValue(tenantReadOutp.Name)
	stateData.DisplayName = types.StringValue(tenantReadOutp.DisplayName)
	stateData.Region = types.StringValue(tenantReadOutp.Region)

	tflog.Debug(ctx, "CDO tenant read")
}

func (t *TenantResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError("Cannot update a created tenant", "Please reach out to CDO TAC if you want to change the display name of your tenant.")
}

func (t *TenantResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	response.Diagnostics.AddError("Cannot delete a created tenant", "Please reach out to CDO TAC if you really want to delete a CDO tenant. You can choose to manually remove the tenant from the Terraform state if you want to remove the tenant from your Terraform configuration.")
}

type PreventUpdatePlanModifier struct{}

func (p PreventUpdatePlanModifier) Description(ctx context.Context) string {
	return "Prevents updates to an existing tenant resource."
}

func (p PreventUpdatePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Prevents updates to an existing tenant resource."
}

// PlanModifyString prevents changes to the string attribute
func (p PreventUpdatePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.StateValue.IsNull() && req.StateValue != req.PlanValue {
		// If the value is changing, prevent it and return an error
		resp.Diagnostics.AddError("Cannot update a created tenant", "Please reach out to CDO TAC if you want to change the display name of your tenant.")
	}
}
