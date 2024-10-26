package msp_tenant

import "github.com/hashicorp/terraform-plugin-framework/types"

type TenantResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DisplayName   types.String `tfsdk:"display_name"`
	GeneratedName types.String `tfsdk:"generated_name"`
	Region        types.String `tfsdk:"region"`
}

type TenantDatasourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Region      types.String `tfsdk:"region"`
}
