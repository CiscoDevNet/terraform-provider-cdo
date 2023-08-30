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
  name             = "<replace-with-email-address-if-non-api-only-user-or-with-username-if-api-only-user>"
  is_api_only_user = "<replace-with-boolean-without-quotes-true-or-false>"
  role             = "<ROLE_READ_ONLY|ROLE_ADMIN|ROLE_SUPER_ADMIN|ROLE_DEPLOY_ONLY|ROLE_EDIT_ONLY|ROLE_VPN_SESSIONS_MANAGER>"
}