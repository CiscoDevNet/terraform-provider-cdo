module "example_vpc" {
  source          = "./vpc"
  resource_prefix = "terraform-provider-example"
}

module "example_sec" {
  source                        = "./sec"
  vpc_id                        = module.example_vpc.vpc_id
  subnet_id                     = module.example_vpc.private_subnet_id
  resource_prefix               = "tf-cdo"
  lb_public_subnet_id           = module.example_vpc.public_subnet_1_id
  lb_secondary_public_subnet_id = module.example_vpc.public_subnet_2_id
  hosted_zone_name              = "<your aws hosted zone>" # e.g. labs.cdo.cisco.com
  dns_prefix                    = "<your choice of sub domain>"  # e.g. test => test.labs.cdo.cisco.com
}

output "sec_bootstrap_data" {
  value = module.example_sec.sec_bootstrap_data
  sensitive = true
}

output "cdo_bootstrap_data" {
  value = module.example_sec.cdo_bootstrap_data
  sensitive = true
}

output "sec_fqdn" {
  value = module.example_sec.sec_fqdn
}

output "sec_instance_id" {
  value = module.example_sec.sec_instance_id
}