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

data "cdo_user" "example_user" {
  name = "example-user@cisco.com"
}

output "example_user_uid" {
  value = data.cdo_user.example_user.id
}

output "example_user_is_api_only_user" {
  value = data.cdo_user.example_user.is_api_only_user
}

output "example_user_role" {
  value = data.cdo_user.example_user.role
}