terraform {
  required_providers {
    cdo = {
      source = "CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "<https://www.defenseorchestrator.com|https://www.defenseorchestrator.eu|https://apj.cdo.cisco.com>"
  api_token = file("${path.module}/api_token.txt")
}