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

resource "cdo_ios_device" "my_ios" {
  name               = "<name-of-device>"
  connector_name     = "<name-of-sdc;not-required-if-connector-type-cdg>"
  socket_address     = "<host>:<port>"
  username           = "<username>"
  password           = "<password>"
  ignore_certificate = "<true|false>"
}
