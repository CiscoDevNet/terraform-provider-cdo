package msp_tenant_user_api_token

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MspManagedTenantUserApiTokenResourceModel struct {
	TenantUid types.String `tfsdk:"tenant_uid"`
	UserUid   types.String `tfsdk:"user_uid"`
	ApiToken  types.String `tfsdk:"api_token"` // Additional field
}
