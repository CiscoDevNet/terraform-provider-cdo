// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	cdoClient "github.com/cisco-lockhart/go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/cisco-lockhart/terraform-provider-cdo/internal/device/asa"
	"github.com/cisco-lockhart/terraform-provider-cdo/internal/device/sdc"
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
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token used to authenticate with the Cisco CDO platform",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The base URL used to communicate with the Cisco CDO platform",
				Required:            true,
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

	client := cdoClient.New(baseURL, apiToken)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CdoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
		asa.NewAsaDeviceResource,
	}
}

func (p *CdoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
		sdc.NewSdcDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CdoProvider{
			version: version,
		}
	}
}
