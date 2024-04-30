terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "<https://www.defenseorchestrator.com|https://www.defenseorchestrator.eu|https://apj.cdo.cisco.com|https://aus.cdo.cisco.com|https://in.cdo.cisco.com>"
  api_token = "<replace-with-api-token-generated-from-cdo>"
}

data "cdo_tenant" "current" {
}

output "current_tenant_uid" {
  value = data.cdo_tenant.current.id
}

output "current_tenant_name" {
  value = data.cdo_tenant.current.name
}

output "current_tenant_human_readable_name" {
  value = data.cdo_tenant.current.human_readable_name
}

output "current_tenant_subscription_type" {
  value = data.cdo_tenant.current.subscription_type
}