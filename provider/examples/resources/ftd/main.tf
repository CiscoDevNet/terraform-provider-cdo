variable "ftd_name" {}

module "vpc" {
  source = "github.com/CiscoDevNet/cdo-automation//modules/aws_vpc"
}

module "bastion" {
  source    = "github.com/CiscoDevNet/cdo-automation//modules/bastion"
  vpc_id    = module.vpc.vpc_id
  subnet_id = module.vpc.public_subnet_id
}

module "ftdv_in_cdo" {
  source              = "github.com/CiscoDevNet/cdo-automation//modules/cdo/ftd"
  bastion_ip          = module.bastion.bastion_ip
  bastion_private_key = module.bastion.bastion_private_key
  bastion_sg          = module.bastion.bastion_sg
  vpc_id              = module.vpc.vpc_id
  public_subnet_id    = module.vpc.public_subnet_id
  private_subnet_id   = module.vpc.private_subnet_id
  ftd_name            = var.ftd_name
  cdo_api_token       = file("${path.module}/api_token.txt")
}