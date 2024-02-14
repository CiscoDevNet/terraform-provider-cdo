terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "<https://www.defenseorchestrator.com|https://www.defenseorchestrator.eu|https://apj.cdo.cisco.com>"
  api_token = "<replace-with-api-token-generated-from-cdo>"
}

resource "cdo_tenant_settings" "tenant_settings" {
  deny_cisco_support_access_to_tenant_enabled = false
}

output "example_tenant_id" {
  value = cdo_tenant_settings.tenant_settings.id
}

output "example_change_request_support_enabled" {
  value = cdo_tenant_settings.tenant_settings.change_request_support_enabled
}

output "example_auto_accept_device_changes_enabled" {
  value = cdo_tenant_settings.tenant_settings.auto_accept_device_changes_enabled
}

output "example_web_analytics_enabled" {
  value = cdo_tenant_settings.tenant_settings.web_analytics_enabled
}

output "example_scheduled_deployments_enabled" {
  value = cdo_tenant_settings.tenant_settings.scheduled_deployments_enabled
}

output "example_deny_cisco_support_access_to_tenant_enabled" {
  value = cdo_tenant_settings.tenant_settings.deny_cisco_support_access_to_tenant_enabled
}

output "example_multi_cloud_defense_enabled" {
  value = cdo_tenant_settings.tenant_settings.multi_cloud_defense_enabled
}

output "example_auto_discover_on_prem_fmcs" {
  value = cdo_tenant_settings.tenant_settings.auto_discover_on_prem_fmcs_enabled
}

output "example_conflict_detection_interval" {
  value = cdo_tenant_settings.tenant_settings.conflict_detection_interval
}
