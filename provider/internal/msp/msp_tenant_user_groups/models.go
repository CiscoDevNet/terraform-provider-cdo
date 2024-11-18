package msp_tenant_user_groups

import "github.com/hashicorp/terraform-plugin-framework/types"

type MspManagedTenantUserGroupsResourceModel struct {
	TenantUid  types.String `tfsdk:"tenant_uid"`
	UserGroups []UserGroup  `tfsdk:"user_groups"`
}

type UserGroup struct {
	Id              types.String `tfsdk:"id"`
	GroupIdentifier types.String `tfsdk:"group_identifier"`
	IssuerUrl       types.String `tfsdk:"issuer_url"`
	Name            types.String `tfsdk:"name"`
	Role            types.String `tfsdk:"role"`
	Notes           types.String `tfsdk:"notes"`
}
