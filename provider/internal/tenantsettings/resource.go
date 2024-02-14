package tenantsettings

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func NewTenantSettingsResource() resource.Resource {
	return &TenantSettingsResource{}
}

type TenantSettingsResource struct {
	client *cdoClient.Client
}

func (*TenantSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_tenant_settings"
}

func (*TenantSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Tenant settings data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the tenant",
				Computed:            true,
			},

			"change_request_support_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether change request support is enabled for the tenant",
				Optional:            true,
				Computed:            true,
			},

			"auto_accept_device_changes_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether auto accept device changes is enabled for the tenant",
				Optional:            true,
				Computed:            true,
			},

			"web_analytics_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether web analytics is enabled for the tenant",
				Optional:            true,
				Computed:            true,
			},

			"scheduled_deployments_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether scheduled deployments is enabled for the tenant",
				Optional:            true,
				Computed:            true,
			},

			"deny_cisco_support_access_to_tenant_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether denying cisco support engineers access to the tenant is enabled",
				Optional:            true,
				Computed:            true,
			},

			"multi_cloud_defense_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether multi cloud defense is enabled for the tenant",
				Optional:            true,
				Computed:            true,
			},

			"auto_discover_on_prem_fmcs_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether change request support is enabled for the tenant",
				Optional:            true,
				Computed:            true,
			},

			"conflict_detection_interval": schema.StringAttribute{
				MarkdownDescription: "The interval used by CDO to detect conflicts on devices",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					validators.NewConflictDetectionIntervalValidator(),
				},
			},
		},
	}
}

func (resource *TenantSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (resource *TenantSettingsResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	handleUpdate(ctx, &res.Diagnostics, &req.Plan, &res.State, resource.client)
}

func (resource *TenantSettingsResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	settings, err := resource.client.ReadTenantSettings(ctx)
	if err != nil {
		res.Diagnostics.AddError("unabled to read tenant settings", err.Error())
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, tenantSettingsDataSourceModelFrom(*settings))...)
}

func (resource *TenantSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	handleUpdate(ctx, &res.Diagnostics, &req.Plan, &res.State, resource.client)
}

func (*TenantSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	res.State.RemoveResource(ctx)
}

func handleUpdate(ctx context.Context, diagnostics *diag.Diagnostics, plan *tfsdk.Plan, state *tfsdk.State, client *cdoClient.Client) {
	var dataModel tenantSettingsDataModel

	diagnostics.Append(plan.Get(ctx, &dataModel)...)
	if diagnostics.HasError() {
		return
	}

	settings, err := client.UpdateTenantSettings(ctx, dataModel.UpdateTenantSettingsInput())
	if err != nil {
		diagnostics.AddError("unable to update tenant settings", err.Error())
		return
	}

	diagnostics.Append(state.Set(ctx, tenantSettingsDataSourceModelFrom(*settings))...)
}
