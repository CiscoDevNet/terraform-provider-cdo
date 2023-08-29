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

resource "cdo_user" "new_user" {
  name             = "jonny.bairstow@ecb.example.com"
  is_api_only_user = false
  role             = "ROLE_ADMIN"
}