terraform {
  required_providers {
    cdo = {
      source = "hashicorp/cisco-lockhart/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://staging.dev.lockhart.io"
  api_token = "<FILL_ME>"
}

data "cdo_asa_device" "my_asa" {
  id = "b66fcd42-f12b-497f-b39e-51fc7d7b8687"
}

output "asa_sdc_type" {
  value = data.cdo_asa_device.my_asa.sdc_type
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
