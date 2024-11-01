data "cdo_msp_managed_tenant" "tenant" {
  name             = "CDO_test-tenant-name"
}

resource "cdo_msp_managed_tenant_users" "example" {
  tenant_uid = data.cdo_msp_managed_tenant.tenant.id
  users = [
    {
      username = "username@example.com",
      roles = ["ROLE_SUPER_ADMIN"]
      api_only_user = false
    },
    {
      username = "username2@example.com",
      roles = ["ROLE_ADMIN"]
      api_only_user = false
    },
    {
      username = "api-only-user",
      roles = ["ROLE_SUPER_ADMIN"]
      api_only_user = true
    }
  ]
}

resource "cdo_msp_managed_tenant_user_api_token" "user_token" {
  tenant_uid = data.cdo_msp_managed_tenant.tenant.id
  user_uid = cdo_msp_managed_tenant_users.example.users[2].id
}

output "api_token" {
  value = cdo_msp_managed_tenant_user_api_token.user_token.api_token
  sensitive = true
}