terraform {
  required_providers {
    cdo = {
      source = "hashicorp/cisco-lockhart/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://ci.dev.lockhart.io"
  api_token = "<FILL_ME>"
}

data "cdo_ios_device" "my_asa" {
  id = "99844204-f604-4acf-b702-d2bdccfabd51"
}