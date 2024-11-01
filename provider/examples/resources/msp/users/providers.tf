terraform {
  required_providers {
    cdo = {
      source = "CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://staging.dev.lockhart.io"
  api_token = file("${path.module}/api_token.txt")
}
