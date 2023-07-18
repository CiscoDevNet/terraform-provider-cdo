terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/cisco-lockhart/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://staging.dev.lockhart.io"
  api_token = "<FILL_ME>"
}

resource "cdo_asa_device" "my_asa" {
  name     = "my_asa"
  sdc_type = "CDG"
  sdc_name = "<FILL_ME>"
  ipv4     = "<FILL_ME>"
  username = "<FILL_ME>"
  password = "<FILL_ME>"
}
