terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/cisco-lockhart/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "<replace-with-cdo-base-url>"
  api_token = "<replace-with-api-token-generated-from-cdo>"
}

data "cdo_asa_device" "my_asa" {
  name = "<enter-device-name>"
}

output "asa_connector_type" {
  value = data.cdo_asa_device.my_asa.connector_type
}
output "asa_sdc_name" {
  value = data.cdo_asa_device.my_asa.sdc_name
}
output "asa_name" {
  value = data.cdo_asa_device.my_asa.name
}
output "asa_socket_address" {
  value = data.cdo_asa_device.my_asa.socket_address
}
output "asa_host" {
  value = data.cdo_asa_device.my_asa.host
}
output "asa_port" {
  value = data.cdo_asa_device.my_asa.port
}
output "asa_ignore_certificate" {
  value = data.cdo_asa_device.my_asa.ignore_certificate
}
