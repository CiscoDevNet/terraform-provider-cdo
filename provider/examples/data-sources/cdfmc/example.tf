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

data "cdo_cdfmc" "current" {
}

output "cdfmc_hostname" {
  value = data.cdo_cdfmc.current.hostname
}

output "cdfmc_software_version" {
  value = data.cdo_cdfmc.current.software_version
}

output "cdfmc_uid" {
  value = data.cdo_cdfmc.current.id
}

output "cdfmc_domain_uuid" {
  value = data.cdo_cdfmc.current.domain_uuid
}