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

resource "cdo_ios_device" "my_ios" {
  name     = "my_ios"
  sdc_type = "SDC"
  sdc_uid  = "<FILL_ME>"
  ipv4     = "<FILL_ME>"
  username = "<FILL_ME>"
  password = "<FILL_ME>"
}
