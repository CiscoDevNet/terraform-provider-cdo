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
  name               = "my_ios"
  connector_type     = "SDC"
  sdc_name           = "<FILL_ME>"
  socket_address     = "<FILL_ME>"
  username           = "<FILL_ME>"
  password           = "<FILL_ME>"
  ignore_certificate = "<FILL_ME>"
}
