terraform {
  required_providers {
    cdo = {
      source = "CiscoDevnet/cdo"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.20.1"
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
}

resource "cdo_ftd_device_onboarding" "test" {
  ftd_uid = cdo_ftd_device.test.id
}