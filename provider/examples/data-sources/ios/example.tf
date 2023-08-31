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

data "cdo_ios_device" "my_ios" {
  name = "<name-of-device>"
}
output "ios_sdc_name" {
  value = data.cdo_ios_device.my_ios.sdc_name
}
output "ios_name" {
  value = data.cdo_ios_device.my_ios.name
}
output "ios_socket_address" {
  value = data.cdo_ios_device.my_ios.socket_address
}
output "ios_host" {
  value = data.cdo_ios_device.my_ios.host
}
output "ios_port" {
  value = data.cdo_ios_device.my_ios.port
}
output "ios_ignore_certificate" {
  value = data.cdo_ios_device.my_ios.ignore_certificate
}