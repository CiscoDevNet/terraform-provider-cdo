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

resource "cdo_asa_device" "my_asa" {
  name               = "<replace-with-name-of-asa>"
  connector_type     = "<CDG|SDC>"
  socket_address     = "<host:port>"
  username           = "<replace-with-username>"
  password           = "<replace-with-password>"
  ignore_certificate = "<replace-with-boolean-without-quotes-true-or-false>"
}
