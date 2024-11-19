data "cdo_msp_managed_tenant" "tenant" {
  name             = "CDO_tenant-name"
}

resource "cdo_msp_managed_tenant_user_groups" "example" {
  tenant_uid = data.cdo_msp_managed_tenant.tenant.id
  user_groups = [
    {
      group_identifier = "customer-managers"
      issuer_url = "https://www.customer-idp.com"
      name =  "customer-managers"
      notes = "Managers in customer's organization"
      role = "ROLE_READ_ONLY"
    },
    {
      group_identifier = "msp-managers"
      issuer_url = "https://www.msp-idp.com"
      name =  "msp-managers"
      notes = "Managers in MSP organization"
      role = "ROLE_READ_ONLY"
    },
    {
      group_identifier = "msp-developers"
      issuer_url = "https://www.msp-idp.com"
      name =  "msp-developers"
      # notes is an optional field, skipped
      role = "ROLE_SUPER_ADMIN"
    }
  ]
}
