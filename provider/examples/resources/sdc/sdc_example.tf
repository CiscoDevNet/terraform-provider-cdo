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

resource "cdo_sdc" "example" {
  name = "tf-sdc-1"
}

output "sdc_name" {
  value = cdo_sdc.example.name
}

output "sdc_uid" {
  value = cdo_sdc.example.id
}

output "sdc_bootstrap_data" {
  value     = cdo_sdc.example.bootstrap_data
  sensitive = true
}