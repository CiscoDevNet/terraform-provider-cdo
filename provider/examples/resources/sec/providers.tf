terraform {
  required_providers {
    cdo = {
      source = "CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://ci.dev.lockhart.io"
  api_token = file("${path.module}/api_token.txt")
}