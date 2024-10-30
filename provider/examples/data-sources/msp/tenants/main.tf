data "cdo_msp_managed_tenant" "tenant" {
  name             = "CDO_isaks-birthday-surprise__skfh2r"
}

output "tenant_display_name" {
  value = data.cdo_msp_managed_tenant.tenant.display_name
}

output "tenant_region" {
  value = data.cdo_msp_managed_tenant.tenant.region
}