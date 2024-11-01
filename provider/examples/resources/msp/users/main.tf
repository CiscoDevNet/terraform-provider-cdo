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
    }
  ]
}