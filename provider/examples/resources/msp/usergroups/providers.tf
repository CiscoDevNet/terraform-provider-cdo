terraform {
  required_providers {
    cdo = {
      source = "CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "https://aus.manage.security.cisco.com"
  api_token = file("${path.module}/api_token.txt")
}
