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
  name               = "burak-crush-lavda"
  connector_type     = "CDG"
  socket_address     = "3.8.235.174:443"
  username           = "lockhart"
  password           = "BlueSkittles123!!"
  ignore_certificate = false
}
