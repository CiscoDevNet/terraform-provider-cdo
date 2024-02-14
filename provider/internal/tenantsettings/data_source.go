package tenantsettings

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NewTenantSettingsDataSource() datasource.DataSource {
	return &TenantSettingsDataSource{}
}

type TenantSettingsDataSource struct {
	client *cdoClient.Client
}

func (*TenantSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_tenant_settings"
}

func (*TenantSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Tenant-wide settings",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the tenant",
				Computed:            true,
			},

			"change_request_support_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether change request support is enabled for the tenant",
				Computed:            true,
			},

			"auto_accept_device_changes_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether auto accept device changes is enabled for the tenant",
				Computed:            true,
			},

			"web_analytics_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether web analytics is enabled for the tenant",
				Computed:            true,
			},

			"scheduled_deployments_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether scheduled deployments is enabled for the tenant",
				Computed:            true,
			},

			"deny_cisco_support_access_to_tenant_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether denying cisco support engineers access to the tenant is enabled",
				Computed:            true,
			},

			"multi_cloud_defense_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether multi cloud defense is enabled for the tenant",
				Computed:            true,
			},

			"auto_discover_on_prem_fmcs_enabled": schema.BoolAttribute{
				MarkdownDescription: "This attribute indicates whether change request support is enabled for the tenant",
				Computed:            true,
			},

			"conflict_detection_interval": schema.StringAttribute{
				MarkdownDescription: "The interval used by CDO to detect conflicts on devices",
				Computed:            true,
			},
		},
	}
}

func (ds *TenantSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

	ds.client = client
}

func (dataSource *TenantSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	settings, err := dataSource.client.ReadTenantSettings(ctx)
	if err != nil {
		res.Diagnostics.AddError("unabled to read tenant settings", err.Error())
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, tenantSettingsDataSourceModelFrom(*settings))...)
}
