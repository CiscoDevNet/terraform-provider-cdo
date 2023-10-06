terraform {
  required_providers {
    cdo = {
      source = "CiscoDevnet/cdo"
    }
  }
}

data "aws_region" "current" {}


output "aws_region" {
  value = data.aws_region.current
}