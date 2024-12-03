resource "cdo_msp_managed_tenant" "tenant" {
  name             = "test-tenant-name"
  display_name = "Display name for tenant"
}

resource "cdo_msp_managed_tenant" "existing_tenant" {
  api_token = "existing-tenant-api-token"
}