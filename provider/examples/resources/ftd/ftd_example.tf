terraform {
  required_providers {
    cdo = {
      source = "hashicorp.com/CiscoDevnet/cdo"
    }
  }
}

provider "cdo" {
  base_url  = "<FILL_ME>"
  api_token = "<FILL_ME>"
}

resource "cdo_ftd_device" "test" {
  name               = "test-weilue-ftd-9"
  access_policy_name = "Default Access Control Policy"
  performance_tier   = "FTDv5"
  virtual            = true
  licenses           = ["BASE"]
}