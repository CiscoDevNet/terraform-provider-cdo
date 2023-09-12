terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27.0"
    }
  }
}

provider "cdo" {
  base_url  = "<FILL_ME>"
  api_token = "<FILL_ME>"
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