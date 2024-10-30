data "cdo_msp_managed_tenant" "tenant" {
  name             = "CDO_tenant-name"
}

output "tenant_display_name" {
  value = data.cdo_msp_managed_tenant.tenant.display_name
}

output "tenant_region" {
  value = data.cdo_msp_managed_tenant.tenant.region
}