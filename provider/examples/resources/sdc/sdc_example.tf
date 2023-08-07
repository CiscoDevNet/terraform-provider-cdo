terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevNet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://staging.dev.lockhart.io"
  api_token = "<FILL_ME>"
}

resource "cdo_sdc" "example" {
  name = "tf-sdc-1"
}

output "sdc_name" {
  value = cdo_sdc.example.name
}

output "sdc_uid" {
  value = cdo_sdc.example.id
}

output "sdc_bootstrap_data" {
  value = cdo_sdc.example.bootstrap_data
}