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


resource "cdo_user" "new_api_only_user" {
  name             = "api_user@kunji.com"
  is_api_only_user = false
  role             = "ROLE_ADMIN"
}

resource "cdo_api_token" "new_api_only_user_api_token" {
    username = cdo_user.new_api_only_user.generated_username
}

output "api_only_user_api_token_value" {
    value = cdo_api_token.new_api_only_user_api_token.api_token
    sensitive = true
}