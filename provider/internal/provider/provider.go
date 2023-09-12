// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/device/ftd/ftdonboarding"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/device/ftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/tenant"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/user"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/user_api_token"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/device/ios"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/device/asa"
)

var _ provider.Provider = &CdoProvider{}

type CdoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CdoProviderModel describes the provider data model.
type CdoProviderModel struct {
	ApiToken types.String `tfsdk:"api_token"`
	BaseURL  types.String `tfsdk:"base_url"`
}

func (p *CdoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cdo"
	resp.Version = p.version
}

func (p *CdoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use the Cisco Defense Orchestrator (CDO) provider to onboard and manage the many devices and other resources supported by CDO. You must configure the provider with the proper credentials and region before you can use it.",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token used to authenticate with CDO. [See here](https://docs.defenseorchestrator.com/c_api-tokens.html#!t-generatean-api-token.html) to learn how to generate an API token.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					validators.OneOfRoles("ROLE_SUPER_ADMIN", "ROLE_ADMIN"),
				},
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base CDO URL. This is the URL you enter when logging into your CDO account.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("https://www.defenseorchestrator.com", "https://www.defenseorchestrator.eu", "https://apj.cdo.cisco.com", "https://staging.dev.lockhart.io", "https://ci.dev.lockhart.io", "https://scale.dev.lockhart.io", "http://localhost:9000"),
				},
			},
		},
	}
}

func (p *CdoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CdoProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown Cisco CDO API Token",
			"The provider cannot create the Cisco CDO API client as there is an unknown configuration value for the Cisco CDO API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CISCO_CDO_API_TOKEN environment variable.",
		)
	}

	if data.BaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Unknown Cisco CDO Base URL",
			"The provider cannot create the Cisco CDO API client as there is an unknown configuration value for the Cisco CDO Base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CISCO_CDO_BASE_URL environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiToken := os.Getenv("CISCO_CDO_API_TOKEN")
	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Cisco CDO API Token",
			"The provider cannot create the Cisco CDO API client as there is a missing or empty value for the Cisco CDO API token. "+
				"Set the API token value in the configuration or use the CISCO_CDO_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	baseURL := os.Getenv("CISCO_CDO_BASE_URL")
	if !data.BaseURL.IsNull() {
		baseURL = data.BaseURL.ValueString()
	}
	if baseURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Missing Cisco CDO Base URL",
			"The provider cannot create the Cisco CDO API client as there is a missing or empty value for the Cisco CDO base URL. "+
				"Set the API token value in the configuration or use the CISCO_CDO_BASE_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := cdoClient.New(baseURL, apiToken)
	if err != nil {
		resp.Diagnostics.AddError("Error while trying to create CDO client", fmt.Sprintf("cause=%s", err.Error()))
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CdoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		connector.NewResource,
		asa.NewAsaDeviceResource,
		ios.NewIosDeviceResource,
		ftd.NewResource,
		user.NewResource,
		user_api_token.NewResource,
		ftdonboarding.NewResource,
	}
}

func (p *CdoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		connector.NewDataSource,
		asa.NewAsaDataSource,
		ios.NewIosDataSource,
		user.NewDataSource,
		tenant.NewDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CdoProvider{
			version: version,
		}
	}
}
