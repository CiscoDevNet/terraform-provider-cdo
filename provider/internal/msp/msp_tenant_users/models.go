package msp_tenant_users

import "github.com/hashicorp/terraform-plugin-framework/types"

type MspManagedTenantUsersResourceModel struct {
	TenantUid types.String `tfsdk:"tenant_uid"`
	Users     []User       `tfsdk:"users"`
}

type User struct {
	Username    types.String `tfsdk:"username"`
	Role        types.String `tfsdk:"role"`
	ApiOnlyUser types.Bool   `tfsdk:"api_only_user"`
}
