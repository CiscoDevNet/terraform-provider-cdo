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

func (*TenantResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_msp_managed_tenant"
}

func (*TenantResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Provides an MSP managed tenant resource. This allows MSP managed tenants to be created. Note: deleting this resource removes the created tenant from the MSP portal by disassociating the tenant from the MSP portal, but the tenant will continue to exist. To completely delete a tenant, please contact Cisco TAC.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the tenant",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the tenant. This should be specified only if a new tenant is being created, and should not be provided if an existing tenant is being added to the MSP protal (i.e., the `api_token` attribute is specified).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					PreventUpdatePlanModifier{}, // Prevent updates to name
				},
				Validators: []validator.String{
					validators.NewMspManagedTenantNameValidator(),
				},
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the tenant. If no display name is specified, the display name will be set to the tenant name. This should be specified only if a new tenant is being created, and should not be provided if an existing tenant is being added to the MSP protal (i.e., the `api_token` attribute is specified).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					PreventUpdatePlanModifier{}, // Prevent updates to name
				},
				Computed: true,
			},
			"generated_name": schema.StringAttribute{
				MarkdownDescription: "Actual name of the tenant returned by the API. This auto-generated name will differ from the name entered by the customer.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "CDO region in which the tenant is created. This is the same region as the region of the MSP portal.",
				Computed:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "API token for an API-only user with super-admin privileges on the tenant. This should be specified only when adding an existing tenant to the MSP portal, and should not be provided if a new tenant is being created (i.e., the `name` and/or `display_name` attributes are specified).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					PreventUpdatePlanModifier{}, // Prevent updates to api token
				},
				Sensitive: true,
			},
		},
	}
}

func (t *TenantResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

	t.client = client
}

func (t *TenantResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating a CDO tenant/Adding an existing tenant using API token to the MSP portal...")

	// 1. Read plan data into planData
	var planData TenantResourceModel

	response.Diagnostics.Append(request.Plan.Get(ctx, &planData)...)

	if response.Diagnostics.HasError() {
		return
	}

	var createOut *tenants.MspTenantOutput
	var err *tenants.CreateError
	if !planData.ApiToken.IsNull() {
		// add tenant to MSP portal
		tflog.Debug(ctx, "Adding existing tenant using API token to MSP portal")
		createOut, err = t.client.AddExistingTenantToMspPortalUsingApiToken(ctx, tenants.MspAddExistingTenantInput{ApiToken: planData.ApiToken.ValueString()})
	} else {
		tflog.Debug(ctx, "Creating new tenant and adding it to the MSP portal")
		// 2. use plan data to create tenant and fill up rest of the model
		createOut, err = t.client.CreateTenantUsingMspPortal(ctx, tenants.MspCreateTenantInput{
			Name:        planData.Name.ValueString(),
			DisplayName: planData.DisplayName.ValueString(),
		})
	}

	if err != nil || createOut == nil {
		response.Diagnostics.AddError("failed to create CDO Tenant", err.Error())
		return
	}

	planData.Id = types.StringValue(createOut.Uid)
	// when a new tenant is created, the name is auto-generated, do not set it to planData.Name
	if planData.Name.IsNull() || planData.Name.IsUnknown() {
		planData.Name = types.StringValue(createOut.Name)
	}
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
	tenantReadOutp, err := t.client.ReadMspManagedTenantByUid(ctx, readInp)

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
	tflog.Debug(ctx, "'Deleting' a CDO tenant by removing it from the MSP portal...")
	var stateData *TenantResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)
	if response.Diagnostics.HasError() {
		return
	}

	deleteInp := tenants.DeleteByUidInput{
		Uid: stateData.Id.ValueString(),
	}
	_, err := t.client.DeleteMspManagedTenantByUid(ctx, deleteInp)
	if err != nil {
		response.Diagnostics.AddError("failed to delete tenant from MSP portal", err.Error())
	}
}

type PreventUpdatePlanModifier struct{}

func (p PreventUpdatePlanModifier) Description(ctx context.Context) string {
	return "Prevents updates to an existing tenant resource."
}

func (p PreventUpdatePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Prevents updates to an existing tenant resource."
}

func (p PreventUpdatePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.StateValue.IsNull() && req.StateValue != req.PlanValue {
		// If the value is changing, prevent it and return an error
		resp.Diagnostics.AddError("Cannot update a created tenant", "Please reach out to CDO TAC if you want to change the display name of your tenant.")
	}
}
