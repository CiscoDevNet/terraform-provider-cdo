module "example_vpc" {
  source = "./aws_vpc"
}

module "example_sdc" {
  source    = "./sdc-aws"
  vpc_id    = module.example_vpc.vpc_id
  subnet_id = module.example_vpc.private_subnet_id
}

output "sdc_aws_region" {
  value = module.example_sdc.aws_region
}