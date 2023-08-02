terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/cisco-lockhart/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "<https://www.defenseorchestrator.com|https://www.defenseorchestrator.eu|https://apj.cdo.cisco.com>"
  api_token = "<replace-with-api-token-generated-from-cdo>"
}

resource "cdo_asa_device" "my_asa" {
  name               = "<name-of-device>"
  connector_type     = "<CDG|SDC>"
  socket_address     = "<host>:<port>"
  username           = "<username>"
  password           = "<password>"
  ignore_certificate = "<true|false>"
}
