terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.20.1"
    }
  }
}

provider "cdo" {
  base_url  = "https://ci.dev.lockhart.io"
  api_token = file("${path.module}/api_token.txt")
}

resource "cdo_ftd_device" "test" {
  name               = "test-wl-ftd"
  access_policy_name = "Default Access Control Policy"
  performance_tier   = "FTDv5"
  virtual            = true
  licenses           = ["BASE"]
  labels             = ["24123", "204", "22", "24", "d", "1", "c", "2", "b", "a", "zzzz", "aaaa", "dddd", "bbbb", "2222", "11111"]
}
