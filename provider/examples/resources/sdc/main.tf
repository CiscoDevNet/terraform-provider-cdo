module "vpc" {
  source = "github.com/CiscoDevNet/cdo-automation//modules/aws_vpc"
}

module "cdo_sdc" {
  source    = "github.com/CiscoDevNet/cdo-automation//modules/cdo/sdc-aws"
  vpc_id    = module.vpc.vpc_id
  subnet_id = module.vpc.private_subnet_id
}

resource "cdo_sdc_onboarding" "sdc" {
  name = module.cdo_sdc.sdc_name
}

output "sdc_name" {
  value = module.cdo_sdc.sdc_name
}
