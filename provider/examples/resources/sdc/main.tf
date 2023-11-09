module "example_vpc" {
  source = "./vpc"
  resource_prefix = "terraform-provider-example"
}

module "example_sdc" {
  source    = "./sdc"
  vpc_id    = module.example_vpc.vpc_id
  subnet_id = module.example_vpc.private_subnet_id
  resource_prefix = "terraform-provider-example"
}

output "sdc_name" {
  value = module.example_sdc.sdc_name
}